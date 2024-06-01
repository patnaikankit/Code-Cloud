// context to manage terminal commands and its functionalities

import React, { createContext, useContext, useState } from "react";

interface CommandType {
    command: string;
    oldDir: string;
    dir: string;
    out: string;
}

interface TerminalContextProps {
    routes: string[];
    setRoutes: (routes: string[]) => void;
    output: CommandType[];
    setOutput: (output: CommandType, clear: boolean) => void;
    activeCommand: string;
    setActiveCommand: (command: string) => void;
}

const TerminalContext = createContext<TerminalContextProps | undefined>(undefined);

export const useTerminal = () => {
    const context = useContext(TerminalContext);

    if (!context) {
        throw new Error('useTerminal must be within a TerminalProvider');
    }

    return context;
}

interface TerminalProviderProps {
    children: React.ReactNode;
}

export const TerminalProvider: React.FC<TerminalProviderProps> = ({ children }) => {
    // to hold the current path
    const [routes, setRoutes] = useState<string[]>(['app']);
    // to hold all the terminal commands
    const [output, setOutput] = useState<CommandType[]>([]);
    // to hold the current active command
    const [activeCommand, setActiveCommand] = useState<string>('');

    // updates the state with new routes
    const setRoutesFn = (routes: string[]) => {
        setRoutes([...routes]);
    }

    // update the state with new command output
    const setOutputFn = (output: CommandType, empty = false) => {
        setActiveCommand('');
        if (empty) {
            setOutput([]);
            return;
        }
        setOutput((e) => [...e, output]);
        setRoutesFn(
            output.dir
                .replace('\r\n', '')
                .split('\\')
                .filter((e) => e !== '')
        );
    }

    return (
        <TerminalContext.Provider
            value={{
                routes,
                setRoutes: setRoutesFn,
                output,
                setOutput: setOutputFn,
                activeCommand,
                setActiveCommand
            }}
        >
            {children}
        </TerminalContext.Provider>
    );
}
