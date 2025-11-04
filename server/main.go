package main

import (
	"context"
	"os"

	"go.lsp.dev/jsonrpc2"
	"go.lsp.dev/protocol"
	"go.uber.org/zap"
)

type handler struct {
	protocol.Server
	logger    *zap.Logger
	documents map[protocol.DocumentURI]string
}

// All functions below have been added as a sort of interface by the protocol.Server
// and all that we are doing is adding the implementation of the behaviour that we want from
// them when they are called by the client.
func (h *handler) Initialize(ctx context.Context, params *protocol.InitializeParams) (*protocol.InitializeResult, error) {
	h.logger.Info("Initialize called")
	return &protocol.InitializeResult{
		Capabilities: protocol.ServerCapabilities{
			// Informs client what capabilities the server has.
			// TextDocumentSync is set to Full which means the server wants to
			// receive the full content of the document
			TextDocumentSync: protocol.TextDocumentSyncOptions{
				OpenClose: true,
				Change:    protocol.TextDocumentSyncKindFull,
			},
			// CompletionProvider indicates that the server provides completion support.
			CompletionProvider: &protocol.CompletionOptions{
				TriggerCharacters: []string{".", ":"},
				ResolveProvider:   false,
			},
			HoverProvider: true,
		},
		// sends the clients information about the LSP name and version.
		ServerInfo: &protocol.ServerInfo{
			Name:    "CC-Tweaked LSP",
			Version: "0.0.1",
		},
	}, nil
}

func (h *handler) Initialized(ctx context.Context, params *protocol.InitializedParams) error {
	h.logger.Info("Initialized called")
	return nil
}

func (h *handler) Shutdown(ctx context.Context) error {
	h.logger.Info("Shutdown called")
	return nil
}

func (h *handler) DidOpen(ctx context.Context, params *protocol.DidOpenTextDocumentParams) error {
	h.logger.Info("Document opened", zap.String("uri", string(params.TextDocument.URI)))
	h.documents[params.TextDocument.URI] = params.TextDocument.Text
	return nil
}

func (h *handler) DidChange(ctx context.Context, params *protocol.DidChangeTextDocumentParams) error {
	h.logger.Info("Document changed", zap.String("uri", string(params.TextDocument.URI)))
	if len(params.ContentChanges) > 0 {
		// Since we use TextDocumentSyncKindFull, we get the full document content
		h.documents[params.TextDocument.URI] = params.ContentChanges[0].Text
	}
	return nil
}

func (h *handler) DidSave(ctx context.Context, params *protocol.DidSaveTextDocumentParams) error {
	h.logger.Info("Document saved", zap.String("uri", string(params.TextDocument.URI)))
	return nil
}

func (h *handler) DidClose(ctx context.Context, params *protocol.DidCloseTextDocumentParams) error {
	h.logger.Info("Document closed", zap.String("uri", string(params.TextDocument.URI)))
	delete(h.documents, params.TextDocument.URI)
	return nil
}

func (h *handler) SetTrace(ctx context.Context, params *protocol.SetTraceParams) error {
	h.logger.Info("SetTrace called", zap.String("value", string(params.Value)))
	return nil
}

func (h *handler) DidChangeWatchedFiles(ctx context.Context, params *protocol.DidChangeWatchedFilesParams) error {
	h.logger.Info("Watched files changed", zap.Int("changes", len(params.Changes)))
	return nil
}

// This struct and functions allows us to use os.Stdin and os.Stdout as a jsonrpc2 stream.
// where we communicate changes.
type stdioConn struct{}

func (s *stdioConn) Read(p []byte) (int, error)  { return os.Stdin.Read(p) }
func (s *stdioConn) Write(p []byte) (int, error) { return os.Stdout.Write(p) }
func (s *stdioConn) Close() error                { return nil }

func main() {
	// Creates a new logger for the server
	logger, _ := zap.NewDevelopment()

	// defer sync to run before closing the application
	defer logger.Sync()

	// We create a new JSON-RPC 2.0 stream over stdin and stdout
	// Where we can both read and write to it.
	stream := jsonrpc2.NewStream(&stdioConn{})

	InitCompletionList()

	// we create a new LSP server in which we pass our handler logger and stream.
	// Note that our handler object here is an instance of a procotol.Server that has a logger attached to it.
	// This means that any method that the server can sent is already set up due to the protocol.Server struct.
	// We have added the logger for our handling of the methods.
	ctx, conn, server := protocol.NewServer(context.Background(), &handler{
		logger:    logger,
		documents: make(map[protocol.DocumentURI]string),
	}, stream, logger)
	if server == nil {
		logger.Fatal("failed to create server")
	}
	// we defer closing the connection until the context is done.
	defer conn.Close()

	<-ctx.Done()
}
