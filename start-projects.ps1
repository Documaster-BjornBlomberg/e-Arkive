# PowerShell script to start both frontend and backend projects

# Start the backend server
Start-Process -NoNewWindow -FilePath "powershell" -ArgumentList "-NoProfile -Command cd backend; go run ."

# Start the frontend server
Start-Process -NoNewWindow -FilePath "powershell" -ArgumentList "-NoProfile -Command npm run dev"