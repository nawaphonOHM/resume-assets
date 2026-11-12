package mongodb

import (
	"context"
	"errors"
	"testing"
)

var (
	errTestCursorDecode = errors.New("decode failure")
	errTestCursorRead   = errors.New("cursor failure")
	errTestCursorClose  = errors.New("close failure")
)

type testAssetCursor struct {
	decodeErr error
	readErr   error
	closeErr  error
	hasNext   bool
	closed    bool
}

func (c *testAssetCursor) Next(context.Context) bool {
	next := c.hasNext
	c.hasNext = false
	return next
}

func (c *testAssetCursor) Decode(any) error { return c.decodeErr }
func (c *testAssetCursor) Err() error       { return c.readErr }

func (c *testAssetCursor) Close(context.Context) error {
	c.closed = true
	return c.closeErr
}

func TestReadAssetRecordsPropagatesDecodeAndCloseErrors(t *testing.T) {
	cursor := &testAssetCursor{hasNext: true, decodeErr: errTestCursorDecode, closeErr: errTestCursorClose}
	records, err := readAssetRecords(t.Context(), cursor)
	if records != nil {
		t.Fatalf("expected no records after decode failure, got %#v", records)
	}
	assertCursorError(t, err, errTestCursorDecode)
	assertCursorError(t, err, errTestCursorClose)
	assertCursorClosed(t, cursor)
}

func TestReadAssetRecordsPropagatesCursorAndCloseErrors(t *testing.T) {
	cursor := &testAssetCursor{readErr: errTestCursorRead, closeErr: errTestCursorClose}
	_, err := readAssetRecords(t.Context(), cursor)
	assertCursorError(t, err, errTestCursorRead)
	assertCursorError(t, err, errTestCursorClose)
	assertCursorClosed(t, cursor)
}

func TestReadAssetRecordsClosesSuccessfulCursor(t *testing.T) {
	cursor := &testAssetCursor{hasNext: true}
	records, err := readAssetRecords(t.Context(), cursor)
	if err != nil {
		t.Fatalf("read records: %v", err)
	}
	if len(records) != 1 {
		t.Fatalf("expected one record, got %d", len(records))
	}
	assertCursorClosed(t, cursor)
}

func assertCursorError(t *testing.T, got, want error) {
	t.Helper()
	if !errors.Is(got, want) {
		t.Fatalf("expected cursor error %v, got %v", want, got)
	}
}

func assertCursorClosed(t *testing.T, cursor *testAssetCursor) {
	t.Helper()
	if !cursor.closed {
		t.Fatal("expected cursor to be closed")
	}
}
