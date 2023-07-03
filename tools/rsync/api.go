// License: GPLv3 Copyright: 2023, Kovid Goyal, <kovid at kovidgoyal.net>

package rsync

import (
	"fmt"
	"io"
	"math"

	"github.com/zeebo/xxh3"

	"kitty/tools/utils"
)

var _ = fmt.Print

const MaxBlockSize int = 256 * 1024

type StrongHashType uint16
type WeakHashType uint16

const (
	XXH3 StrongHashType = iota
)
const (
	Beta WeakHashType = iota
)

type Api struct {
	rsync                                            RSync
	signature                                        []BlockHash
	delta_output                                     io.Writer
	delta_input                                      io.ReadSeeker
	unconsumed_signature_data, unconsumed_delta_data []byte
	expected_input_size_for_signature_generation     int64

	Strong_hash_type StrongHashType
	Weak_hash_type   WeakHashType
}

// internal implementation {{{
func (self *Api) read_signature_header(data []byte) (consumed int, err error) {
	if len(data) < 12 {
		return -1, io.ErrShortBuffer
	}
	if version := bin.Uint32(data); version != 0 {
		return consumed, fmt.Errorf("Invalid version in signature header: %d", version)
	}
	switch strong_hash := StrongHashType(bin.Uint16(data[4:])); strong_hash {
	case XXH3:
		self.Strong_hash_type = strong_hash
		self.rsync.UniqueHasher = xxh3.New()
	default:
		return consumed, fmt.Errorf("Invalid strong_hash in signature header: %d", strong_hash)
	}
	switch weak_hash := WeakHashType(bin.Uint16(data[6:])); weak_hash {
	case Beta:
		self.Weak_hash_type = weak_hash
	default:
		return consumed, fmt.Errorf("Invalid weak_hash in signature header: %d", weak_hash)
	}
	block_size := int(bin.Uint32(data[8:]))
	consumed = 12
	if block_size == 0 {
		return consumed, fmt.Errorf("rsync signature header has zero block size")
	}
	if block_size > MaxBlockSize {
		return consumed, fmt.Errorf("rsync signature header has too large block size %d > %d", block_size, MaxBlockSize)
	}
	self.rsync.BlockSize = block_size
	self.rsync.MaxDataOp = 10 * block_size
	self.signature = make([]BlockHash, 0, 256)
	return
}

func (self *Api) read_signature_blocks(data []byte) (consumed int) {
	hash_size := self.rsync.UniqueHasher.Size()
	block_hash_size := hash_size + 12
	for ; len(data) >= block_hash_size; data = data[block_hash_size:] {
		bl := BlockHash{}
		bl.Unserialize(data[:block_hash_size], hash_size)
		self.signature = append(self.signature, bl)
		consumed += block_hash_size
	}
	return
}

func (self *Api) update_delta(data []byte) (consumed int, err error) {
	op := Operation{}
	for len(data) > 0 {
		n, uerr := op.Unserialize(data)
		if uerr == nil {
			consumed += n
			data = data[n:]
			if err = self.rsync.ApplyDelta(self.delta_output, self.delta_input, op); err != nil {
				return
			}
		} else {
			if n < 0 {
				return consumed, nil
			}
			return consumed, uerr
		}
	}
	return
}

// }}}

// Start applying serialized delta
func (self *Api) StartDelta(delta_output io.Writer, delta_input io.ReadSeeker) {
	self.delta_output = delta_output
	self.delta_input = delta_input
	self.unconsumed_delta_data = nil
}

// Apply a chunk of delta data
func (self *Api) UpdateDelta(data []byte) (err error) {
	if len(self.unconsumed_delta_data) > 0 {
		data = append(self.unconsumed_delta_data, data...)
		self.unconsumed_delta_data = nil
	}
	consumed, err := self.update_delta(data)
	if err != nil {
		return err
	}
	data = data[consumed:]
	if len(data) > 0 {
		self.unconsumed_delta_data = data
	}
	return
}

// Finish applying delta data
func (self *Api) FinishDelta() (err error) {
	if err = self.UpdateDelta([]byte{}); err != nil {
		return err
	}
	if len(self.unconsumed_delta_data) > 0 {
		return fmt.Errorf("There are %d leftover bytes in the delta", len(self.unconsumed_delta_data))
	}
	self.delta_input = nil
	self.delta_output = nil
	self.unconsumed_delta_data = nil
	return
}

// Create a serialized delta based on the previously loaded signature
func (self *Api) CreateDelta(src io.Reader, output_callback func(string) error) (err error) {
	if len(self.signature) == 0 {
		return fmt.Errorf("Cannot call CreateDelta() before loading a signature")
	}
	self.rsync.CreateDelta(src, self.signature, func(op Operation) error {
		return output_callback(op.Serialize())
	})
	return
}

// Create a signature for the data source in src
func (self *Api) CreateSignature(src io.Reader, callback func([]byte) error) (err error) {
	header := make([]byte, 12)
	bin.PutUint16(header[4:], uint16(self.Strong_hash_type))
	bin.PutUint16(header[6:], uint16(self.Weak_hash_type))
	bin.PutUint32(header[8:], uint32(self.rsync.BlockSize))
	if err = callback(header); err != nil {
		return err
	}
	if self.expected_input_size_for_signature_generation > 0 {
		self.signature = make([]BlockHash, 0, self.rsync.BlockHashCount(self.expected_input_size_for_signature_generation))
	} else {
		self.signature = make([]BlockHash, 0, 1024)
	}
	return self.rsync.CreateSignature(src, func(bl BlockHash) error {
		if err = callback(bl.Serialize()); err != nil {
			return err
		}
		self.signature = append(self.signature, bl)
		return nil
	})
}

// Add more external signature data
func (self *Api) AddSignatureData(data []byte) (err error) {
	if len(self.unconsumed_signature_data) > 0 {
		data = append(self.unconsumed_signature_data, data...)
		self.unconsumed_signature_data = nil
	}
	if self.rsync.UniqueHasher == nil {
		consumed, err := self.read_signature_header(data)
		if err != nil {
			if consumed < 0 {
				self.unconsumed_signature_data = data
				return nil
			}
			return err
		}
		data = data[consumed:]
	}
	consumed := self.read_signature_blocks(data)
	data = data[consumed:]
	if len(data) > 0 {
		self.unconsumed_signature_data = data
	}
	return nil
}

// Finish adding external signature data
func (self *Api) FinishSignatureData() (err error) {
	if len(self.unconsumed_signature_data) > 0 {
		return fmt.Errorf("There were %d leftover bytes in the signature data", len(self.unconsumed_signature_data))
	}
	self.unconsumed_signature_data = nil
	if self.rsync.UniqueHasher == nil {
		return fmt.Errorf("No header was found in the signature data")
	}
	return
}

// Use to calculate a delta based on a supplied signature, via AddSignatureData
func NewToCreateDelta() *Api {
	return &Api{}
}

// Use to create a signature and possibly apply a delta
func NewToCreateSignature(expected_input_size int64) (ans *Api, err error) {
	bs := DefaultBlockSize
	sz := utils.Max(0, expected_input_size)
	if sz > 0 {
		bs = int(math.Round(math.Sqrt(float64(sz))))
	}
	ans = &Api{}
	ans.rsync.BlockSize = utils.Min(bs, MaxBlockSize)
	ans.rsync.UniqueHasher = xxh3.New()

	if ans.rsync.UniqueHasher.BlockSize() > 0 && ans.rsync.UniqueHasher.BlockSize() < ans.rsync.BlockSize {
		ans.rsync.BlockSize = (ans.rsync.BlockSize / ans.rsync.UniqueHasher.BlockSize()) * ans.rsync.UniqueHasher.BlockSize()
	}

	ans.rsync.MaxDataOp = ans.rsync.BlockSize * 10
	ans.expected_input_size_for_signature_generation = sz
	return
}