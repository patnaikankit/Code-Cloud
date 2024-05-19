// context to manage websocket connections 

import React, { ReactNode, createContext, useContext, useState } from "react"
import { useTerminal } from "./TerminalUtil"
import { useFiles } from "./FileUtil"

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
    // to hold websocket object 
    const [socket, setSocket] = useState<WebSocket | null>(null)
    const {routes, activeCommand, setOutput} = useTerminal()
    const {setFiles, setFileData} = useFiles()

    // establish new socket connection
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
            if(response.type === 'files'){
                setFiles(
                    response.dir.split('/app/')[1] || '',
                    response.out.split('\n').filter((s: string) => s !== '')
                )
            }
            else if(response.type === 'file'){
                setFileData(
                    response.dir.split('/app/')[1] + '/' + response.isFile || '' + response.isFile,
                    response.out
                )
            }
            else if(response.type === 'command'){
                if(response.error){
                    setOutput({
                        dir: '/' + routes.join('/'),
                        oldDir: response.oldDir,
                        command: response.command,
                        out: response.out + response.error
                    },
                    false
                )}
                else{
                    setOutput({
                        dir: response.dir,
                        oldDir: response.oldDir,
                        command: response.command,
                        out: response.out
                    }, false)
                }
            }
        })
        setSocket(newSocket)
    }

    const sendMessage = (message: string) => {
        if(socket && socket.readyState === WebSocket.OPEN){
            socket.send(message)
        }
    }

    // construct and send a command to get file
    const getFile = (path: string) => {
        const filterPath = path.split('/').filter((e) => e !== '')
        const removeFile = filterPath.slice(0, -1).join('/')
        const fileName = filterPath.pop()
        const cmd = {
            dir: '/app/' + removeFile,
            command: 'cat' + fileName,
            type: 'file',
            isFile: fileName
        }
        if(socket && socket.readyState === WebSocket.OPEN){
            socket.send(JSON.stringify(cmd))
        }
    }

    // construct and send a command to save a file
    const saveFile = (path: string, data: string) => {
        const filterPath = path.split('/').filter((e) => e !== '')
        const removeFile = filterPath.slice(0, -1).join('/')
        const fileName = filterPath.pop()
        const cmd = {
            dir: '/app/' + removeFile,
            command: '',
            type: 'file',
            isFile: fileName,
            isCustom: false,
            data: data || ' '
        }
        if(socket && socket.readyState === WebSocket.OPEN){
            socket.send(JSON.stringify(cmd))
        }
    }

    // execute the active command
    const execCommand =() => {
        if(activeCommand === ''){
            return
        }
        if(activeCommand === 'clear'){
            setOutput({
                dir: '/' + routes.join('/'),
                oldDir: '',
                command: activeCommand.trim(),
                out: ''
            }, true)
            return
        }
        const cmd = {
            dir: '/' + routes.join('/'),
            command: activeCommand.trim(),
            type: 'command',
            isFile: ''
        }
        if(socket && socket.readyState === WebSocket.OPEN){
            socket.send(JSON.stringify(cmd))
        }
    } 

    const contextValue: WebSocketContextProps = {
        socket,
        setSocket: setSocketFn,
        sendMessage,
        getFile,
        saveFile,
        execCommand
    }

    return (
        <WebSocketContext.Provider value={contextValue}>{children}</WebSocketContext.Provider>
    )
}