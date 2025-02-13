package web

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func Respond(ctx context.Context, w http.ResponseWriter, statusCode int, data any) error {
	if err := ctx.Err(); err != nil {
		if errors.Is(err, context.Canceled) {
			return errors.New("client cancelled the request")
		}
	}

	//no content
	if statusCode == http.StatusNoContent {
		//just write the header
		w.WriteHeader(statusCode)
		return nil
	}

	bs, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	w.WriteHeader(statusCode)
	if _, err := w.Write(bs); err != nil {
		return fmt.Errorf("writing response: %w", err)
	}
	return nil
}
