package types

import "time"

type Circuit struct {
	Opens  map[string]Operation
	Closed map[string]Operation
}

func NewCircuit() *Circuit {
	return &Circuit{
		Opens:  map[string]Operation{},
		Closed: map[string]Operation{},
	}
}

func (c *Circuit) Run() {
	for {
		for _, operation := range c.Opens {
			if operation.Locked {
				continue
			}

			err := make(chan error)

			go operation.Exec(err)

			select {
			case err := <-err:
				if err != nil {
					go operation.Wait()
					break
				}
				go c.Close(&operation)
			case <-time.After(time.Duration(operation.Timeout)):
				go operation.Wait()
			}
		}
	}
}

func (c *Circuit) Append(operation *Operation) {
	c.Opens[operation.ID.String()] = *operation
}

func (c *Circuit) Close(operation *Operation) {
	delete(c.Opens, operation.ID.String())
	c.Closed[operation.ID.String()] = *operation
}
