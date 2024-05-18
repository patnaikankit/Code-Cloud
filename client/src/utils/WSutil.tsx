import React, { ReactNode, createContext, useContext, useState } from "react"
import { useTerminal } from "./TerminalUtil"
import { useFiles } from "./FIleUtil"

interface WebSocketContextProps {
    socket: WebSocket | null
    setSocket: (id: string) => void
    sendMessage: (path: string) => void
    getFile: (path: string) => void
    saveFile: (path: string, data: string) => void
    execCommand: () => void 
}

const WebSocketContext = createContext<WebSocketContextProps | undefined>(undefined)

export const useWebsocket = () => {
    const context = useContext(WebSocketContext)

    if(!context){
        throw new Error('useWebscoket should be within a webSocketProvider')
    }
    return context
}

interface WebSocketProviderProps{
    children: ReactNode
}

export const WebSocketProvider: React.FC<WebSocketProviderProps> = ({ children }) => {
    const [socket, setSocket] = useState<WebSocket | null>(null)
    const {routes, activeCommand, setOutput} = useTerminal()
    const {setFile, setFileData} = useFiles()

    const setSocketFn = (id: string) => {
        if(socket){
            socket.close()
        }

        const newSocket = new WebSocket(`ws://${id}.localhost:5000/ws`)

        newSocket.addEventListener('open', (e) => {
            console.log('Websocket connection opened:',e);
        })
        
        newSocket.addEventListener('close', (e) => {
            console.log('Websocket connection closed:',e);
        })

        newSocket.addEventListener('error', (e) => {
            console.log('Websocket connection error:',e);
        })

        newSocket.addEventListener('message', (e) => {
            const response = JSON.parse(e.data)
            if(response === 'files'){
                setFile(
                    response.dir.split('/app/')[1] || '',
                    response.out.split('\n').filter((s: string) => s !== '')
                )
            }
            else if(response === 'file'){
                response.dir.split('/app/')[1] + '/' + response.isFile || '' + response.isFile,
                response.out
            }
        })
    }
}