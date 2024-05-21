import ReactDOM from 'react-dom/client'
import App from './App.tsx'
import './index.css'
import { WebSocketProvider } from './utils/WSutil.tsx'
import { FileProvider } from './utils/FileUtil.tsx'
import { TerminalProvider } from './utils/TerminalUtil.tsx'

ReactDOM.createRoot(document.getElementById('root')!).render(
  <TerminalProvider>
    <FileProvider>
      <WebSocketProvider>
        <App />
      </WebSocketProvider>
    </FileProvider>
  </TerminalProvider>
)
