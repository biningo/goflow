package core

import (
	"testing"
)

func TestJob(t *testing.T) {
	j := Job("job 1")

	add_1_1 := Task("add 1 1", AddOperator(1, 1))
	sleep_2 := Task("sleep 2", SleepOperator(2))
	add_2_4 := Task("add 2 4", AddOperator(2, 4))
	add_3_4 := Task("add 3 4", AddOperator(3, 4))

	j.addTask(add_1_1)
	j.addTask(sleep_2)
	j.addTask(add_2_4)
	j.addTask(add_3_4)

	j.setDownstream(add_1_1, sleep_2)
	j.setDownstream(sleep_2, add_2_4)
	j.setDownstream(add_1_1, add_3_4)

	j.run()

	b := Job("bad job")
	add_2_2 := Task("add 2 2", AddOperator(2, 2))
	add_4_4 := Task("add 4 4", AddOperator(4, 4))
	b.addTask(add_2_2)
	b.addTask(add_4_4)
	b.setDownstream(add_2_2, add_4_4)
	b.setDownstream(add_4_4, add_2_2)
	b.run()
}
