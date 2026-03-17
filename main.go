package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"syscall"
)

type DiskStatus struct {
	Percent float64 `json:"percent"`
	Used    uint64  `json:"used"`
	Free    uint64  `json:"free"`
}

func getDiskApi(w http.ResponseWriter, r *http.Request) {
	var stat syscall.Statfs_t
	syscall.Statfs("/", &stat)
	all := stat.Blocks * uint64(stat.Bsize)
	free := stat.Bfree * uint64(stat.Bsize)
	used := all - free
	percent := (float64(used) / float64(all)) * 100

	status := DiskStatus{
		Percent: percent,
		Used:    used / (1024 * 1024 * 1024),
		Free:    free / (1024 * 1024 * 1024),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

func main() {
	fmt.Println("🚀 SRE Full-Stack Monitor starting on :8081...")

	// 1. 首页路由
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "<h1>SRE GO BACKEND</h1><p>Visit <a href='/dashboard'>/dashboard</a> for the UI.</p>")
	})

	// 2. 数据接口
	http.HandleFunc("/api/disk", getDiskApi)

	// 3. React 仪表盘路由 (核心)
	http.HandleFunc("/dashboard", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `
<!DOCTYPE html>
<html>
<head>
    <title>GO-REACT DASHBOARD</title>
    <script src="https://unpkg.com/react@18/umd/react.development.js"></script>
    <script src="https://unpkg.com/react-dom@18/umd/react-dom.development.js"></script>
    <script src="https://unpkg.com/@babel/standalone/babel.min.js"></script>
    <style>
        body { background: #0a0a0a; color: #00d8ff; font-family: 'Segoe UI', sans-serif; display: flex; justify-content: center; align-items: center; height: 100vh; margin: 0; }
        .card { border: 2px solid #00d8ff; padding: 30px; border-radius: 15px; background: #111; box-shadow: 0 0 30px rgba(0, 216, 255, 0.2); text-align: center; width: 400px; }
        .progress-container { background: #333; height: 25px; border-radius: 12px; margin: 20px 0; overflow: hidden; border: 1px solid #444; }
        .progress-fill { background: linear-gradient(90deg, #00d8ff, #005f7a); height: 100%; transition: width 0.8s cubic-bezier(0.4, 0, 0.2, 1); }
        h1 { font-size: 1.8em; margin-bottom: 10px; text-transform: uppercase; letter-spacing: 2px; }
        .stats { color: #aaa; font-size: 0.9em; }
    </style>
</head>
<body>
    <div id="root"></div>
    <script type="text/babel">
        function App() {
            const [disk, setDisk] = React.useState({ percent: 0, used: 0, free: 0 });

            React.useEffect(() => {
                const fetchData = () => {
                    fetch('/api/disk')
                        .then(res => res.json())
                        .then(data => setDisk(data));
                };
                fetchData();
                const timer = setInterval(fetchData, 3000); // 每3秒自动同步后端数据
                return () => clearInterval(timer);
            }, []);

            return (
                <div className="card">
                    <h1>System Monitor</h1>
                    <div className="progress-container">
                        <div className="progress-fill" style={{ width: disk.percent + "%" }}></div>
                    </div>
                    <div className="stats">
                        <p>DISK USAGE: <strong>{disk.percent.toFixed(1)}%</strong></p>
                        <p>USED: {disk.used} GB / FREE: {disk.free} GB</p>
                    </div>
                    <p style={{marginTop: '20px', fontSize: '0.7em', color: '#555'}}>Status: Real-time via Go API</p>
                </div>
            );
        }
        const root = ReactDOM.createRoot(document.getElementById('root'));
        root.render(<App />);
    </script>
</body>
</html>
		`)
	})

	http.ListenAndServe(":8081", nil)
}