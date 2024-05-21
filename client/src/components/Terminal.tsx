// terminal interface

import React, { useEffect, useRef } from "react";
import { useTerminal } from "../utils/TerminalUtil";
import { useWebsocket } from "../utils/WSutil";

const Terminal = () => {
    const { output, routes, activeCommand, setActiveCommand } = useTerminal()
    const { execCommand } = useWebsocket()
    const terminalRef = useRef<HTMLDivElement>(null)
    const inputRef = useRef<HTMLInputElement>(null)

    // execute command
    const execCmd = (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault()
        execCommand()
    }

    const scrollToBottom = () => {
        if(terminalRef.current){
            terminalRef.current.scrollTop = terminalRef.current.scrollHeight
        }
    }

    useEffect(() => {
        scrollToBottom()
    }, [output])

    return (
        <div
         className="h-64 bg-fStructBackground p-2 overflow-auto"
         id="terminal"
         ref={terminalRef}
         onClick={() => {
            inputRef.current?.focus()
         }}
        >
            <ul>
                {/* list of  previous terinal commands */}
                {output.map((out, i) => {
                    return (
                        <li
                            key={out.command + out.dir + i}
                        >
                            <div className="flex text-sm gap-2 items-center">
                                <span className="text-blue-600 flex">
                                    <pre className="flex">
                                        {out.oldDir
                                            .split('/')
                                            .filter((e) => e !== '')
                                            .join('/')
                                        }
                                    </pre>
                                    <pre></pre>
                                </span>
                                <pre 
                                    key={out.command}
                                    className="bg-transparent outline-none text-gray-500"
                                >
                                    {out.command}
                                </pre>
                            </div>
                            <pre className="text-gray-500 text-xs">{out.out}</pre>
                        </li>
                    );
                })}
            </ul>
            {/* command inputs */}
            <form 
                className="flex text-sm gap-2 items-center"
                onSubmit={execCmd}
                id="terminal-input"
            >
                <span className="text-blue-600 flex">
                    <pre>{routes.join('/')}</pre>
                    <pre></pre>
                </span>
                <pre>
                    <input
                        type="text"
                        ref={inputRef}
                        className="bg-transparent outline-none text-gray-500 w-full"
                        value={activeCommand}
                        onChange={(e) => setActiveCommand(e.target.value)}
                    />
                </pre>
            </form>
        </div>
    )
}

export default Terminal;