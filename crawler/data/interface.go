package data

type DataOperationType int

type DataReader interface {
	Read(type_ DataOperationType, args ...interface{}) (interface{}, error)
}

type DataWriter interface {
	Write(type_ DataOperationType, args ...interface{}) error
}

type DataReadWriter interface {
	DataReader
	DataWriter
}
