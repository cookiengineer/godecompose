package disasm

import "sort"

// BasicBlock represents a straight-line sequence of instructions with a single
// entry point and potentially multiple exit points.
type BasicBlock struct {
	StartAddr    uint64
	EndAddr      uint64
	Instructions []Instruction
	Successors   []*BasicBlock
	Predecessors []*BasicBlock
}

// BuildControlFlowGraph constructs a control flow graph from a linear
// instruction stream. entryPoints specifies known function entry addresses.
func BuildControlFlowGraph(instructions []Instruction, entryPoints []uint64) []*BasicBlock {
	if len(instructions) == 0 {
		return nil
	}

	addrToIndex := make(map[uint64]int)
	for i, inst := range instructions {
		addrToIndex[inst.Address] = i
	}

	leaders := make(map[uint64]bool)
	leaders[instructions[0].Address] = true

	for _, ep := range entryPoints {
		leaders[ep] = true
	}

	for _, inst := range instructions {
		if inst.IsBranch || inst.IsCall {
			if inst.BranchTarget != 0 {
				if _, exists := addrToIndex[inst.BranchTarget]; exists {
					leaders[inst.BranchTarget] = true
				}
			}
		}

		nextAddr := inst.Address + uint64(inst.Size)

		if inst.IsReturn || inst.Opcode == "RET" {
			if nextIdx, ok := addrToIndex[nextAddr]; ok && nextIdx < len(instructions) {
				leaders[nextAddr] = true
			}
			continue
		}

		if inst.Opcode == "JMP" || inst.Opcode == "JMPL" || inst.Opcode == "JMPQ" {
			if nextIdx, ok := addrToIndex[nextAddr]; ok && nextIdx < len(instructions) {
				leaders[nextAddr] = true
			}
			continue
		}

		if inst.IsCall {
			if nextIdx, ok := addrToIndex[nextAddr]; ok && nextIdx < len(instructions) {
				leaders[nextAddr] = true
			}
			continue
		}
	}

	sortedLeaders := make([]uint64, 0, len(leaders))
	for addr := range leaders {
		sortedLeaders = append(sortedLeaders, addr)
	}
	sort.Slice(sortedLeaders, func(i, j int) bool {
		return sortedLeaders[i] < sortedLeaders[j]
	})

	blocks := make([]*BasicBlock, 0, len(sortedLeaders))
	blockByAddr := make(map[uint64]*BasicBlock)

	for i, leaderAddr := range sortedLeaders {
		startIdx := addrToIndex[leaderAddr]
		endIdx := len(instructions)

		if i+1 < len(sortedLeaders) {
			nextLeader := sortedLeaders[i+1]
			if idx, ok := addrToIndex[nextLeader]; ok {
				endIdx = idx
			}
		}

		blockInsts := instructions[startIdx:endIdx]
		block := &BasicBlock{
			StartAddr:    leaderAddr,
			EndAddr:      blockInsts[len(blockInsts)-1].Address + uint64(blockInsts[len(blockInsts)-1].Size),
			Instructions: blockInsts,
		}

		blocks = append(blocks, block)
		blockByAddr[leaderAddr] = block
	}

	for _, block := range blocks {
		lastInst := block.Instructions[len(block.Instructions)-1]

		for _, inst := range block.Instructions {
			if inst.IsCall && inst.BranchTarget != 0 {
				if callee, ok := blockByAddr[inst.BranchTarget]; ok {
					addEdge(block, callee)
				}
			}
			if inst.IsBranch && inst.BranchTarget != 0 && !inst.IsCall {
				if target, ok := blockByAddr[inst.BranchTarget]; ok {
					addEdge(block, target)
				}
			}
		}

		if lastInst.IsReturn {
			continue
		}

		if lastInst.Opcode == "JMP" || lastInst.Opcode == "JMPL" || lastInst.Opcode == "JMPQ" {
			continue
		}

		nextAddr := lastInst.Address + uint64(lastInst.Size)
		if next, ok := blockByAddr[nextAddr]; ok {
			addEdge(block, next)
		}
	}

	return blocks
}

func addEdge(from, to *BasicBlock) {
	from.Successors = append(from.Successors, to)
	to.Predecessors = append(to.Predecessors, from)
}
