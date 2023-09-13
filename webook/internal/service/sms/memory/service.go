package memory

import (
	"context"
	"fmt"
)

type MemoryService struct {
}

func NewMemoryService() *MemoryService {
	return &MemoryService{}
}
func (m *MemoryService) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	fmt.Println(args)
	return nil
}
