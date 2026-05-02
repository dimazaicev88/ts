package base

type WorkerStatus int8

const Disconnected WorkerStatus = 1
const Connected WorkerStatus = 1

func (w WorkerStatus) ToInt8() int8 {
	return int8(w)
}
