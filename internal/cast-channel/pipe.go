package castchannel

type pipeCastMessageReader struct {
	c chan *CastMessage
}

func (t *pipeCastMessageReader) Read(cm *CastMessage) error {
	*cm = *(<-t.c)
	return nil
}

type pipeCastMessageWriter struct {
	c chan *CastMessage
}

func (t *pipeCastMessageWriter) Write(cm *CastMessage) error {
	t.c <- cm
	return nil
}

func Pipe() (*pipeCastMessageReader, *pipeCastMessageWriter) {
	c := make(chan *CastMessage)
	r := &pipeCastMessageReader{c}
	w := &pipeCastMessageWriter{c}

	return r, w
}
