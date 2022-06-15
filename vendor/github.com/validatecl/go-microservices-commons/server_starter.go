package commons

import (
	"net"
	"net/http"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"github.com/oklog/oklog/pkg/group"
)

//CreateServer inicializa un http server
func CreateServer(handler http.Handler, port string, logger log.Logger) group.Group {
	var g group.Group
	{

		{
			httpListener, err := net.Listen("tcp", port)
			if err != nil {
				level.Info(logger).Log("transport", "HTTP", "during", "Listen", "err", err)
				os.Exit(1)
			}

			g.Add(func() error {
				level.Info(logger).Log("transport", "HTTP", "addr", port)
				return http.Serve(httpListener, handler)
			}, func(error) {
				httpListener.Close()
			})
		}
	}

	return g
}
