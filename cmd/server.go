package cmd

import (
	"fmt"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/yechentide/necrack/server"
	"github.com/yechentide/necrack/styles"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start HTTP server for ZIP file processing",
	Long: `Start an HTTP server that accepts ZIP file uploads containing NetEase Minecraft worlds,
decrypts them, and returns the processed files as a ZIP download.

The server provides a single endpoint:
  POST /decrypt - Upload a ZIP file and receive the decrypted version

Example:
  necrack server --port 8080

  # Upload and decrypt a ZIP file using curl:
  curl -X POST -F "zipfile=@world.zip" http://localhost:8080/decrypt -o decrypted.zip`,
	Run: func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetInt("port")
		
		// Setup styled output from centralized styles
		
		// Setup logger
		logger := log.NewWithOptions(nil, log.Options{
			ReportTimestamp: true,
			TimeFormat:      "15:04:05",
			Prefix:          "[necrack-server]",
		})
		log.SetDefault(logger)
		
		http.HandleFunc("/decrypt", server.DecryptHandler)
		
		// Health check endpoint
		http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		})
		
		// Simple landing page
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/" {
				http.NotFound(w, r)
				return
			}
			
			html := `<!DOCTYPE html>
<html>
<head>
    <title>NetEase World Decryption Service</title>
    <style>
        body { font-family: Arial, sans-serif; max-width: 600px; margin: 50px auto; padding: 20px; }
        .upload-form { border: 2px dashed #ccc; padding: 30px; text-align: center; margin: 20px 0; }
        input[type="file"] { margin: 10px 0; }
        button { background: #007cba; color: white; padding: 10px 20px; border: none; border-radius: 5px; cursor: pointer; }
        button:hover { background: #005a87; }
    </style>
</head>
<body>
    <h1>NetEase World Decryption Service</h1>
    <p>Upload a ZIP file containing NetEase Minecraft worlds to decrypt them.</p>
    
    <div class="upload-form">
        <form action="/decrypt" method="post" enctype="multipart/form-data">
            <p>Select a ZIP file to decrypt:</p>
            <input type="file" name="zipfile" accept=".zip" required>
            <br>
            <button type="submit">Upload and Decrypt</button>
        </form>
    </div>
    
    <h3>API Usage:</h3>
    <pre>curl -X POST -F "zipfile=@world.zip" http://localhost:%d/decrypt -o decrypted.zip</pre>
</body>
</html>`
			
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprintf(w, html, port)
		})

		addr := fmt.Sprintf(":%d", port)
		
		// Display startup information with styling
		fmt.Println(styles.ServerHeaderStyle.Render("🌍 NetEase World Decryption Server"))
		fmt.Println()
		fmt.Printf("%s %s\n", 
			styles.InfoStyle.Render("⚡ Server starting on:"), 
			styles.URLStyle.Render(fmt.Sprintf("http://localhost:%d", port)))
		fmt.Printf("%s %s\n", 
			styles.InfoStyle.Render("📤 Upload endpoint:"), 
			styles.URLStyle.Render(fmt.Sprintf("http://localhost:%d/decrypt", port)))
		fmt.Printf("%s %s\n", 
			styles.InfoStyle.Render("💚 Health check:"), 
			styles.URLStyle.Render(fmt.Sprintf("http://localhost:%d/health", port)))
		fmt.Println()
		
		logger.Info("Server starting", "port", port, "addr", addr)
		
		if err := http.ListenAndServe(addr, nil); err != nil {
			logger.Fatal("Server failed to start", "error", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().IntP("port", "p", 8080, "Port to run the server on")
}