package event

type ServerEventHandler struct {
	Handlers map[EEvent]HandleFunc
}

func handlePrintEvent(evt PrintEvent) error {
	_, err := evt.Encode()
	if err != nil {
		return err
	}

	return nil
}
