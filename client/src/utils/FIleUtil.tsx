// context to manage files and folders and their functionalities

import React, { createContext, useContext, useEffect, useState } from "react";

interface FileType {
    type: string;
    path: string;
    data?: string;
}

interface FolderType {
    type: string;
    path: string;
}

interface FileContextProps {
    files: { [key: string]: FileType | FolderType };
    setFiles: (dir: string, tempFiles: string[]) => void;
    setFileData: (dir: string, data: string) => void;
    activeFile: string;
    activeFileData: string;
    removeFile: (dir: string) => void;
    removeFolder: (dir: string) => void;
}

const FileUtils = createContext<FileContextProps | undefined>(undefined);

export const useFiles = () => {
    const context = useContext(FileUtils);

    if (!context) {
        throw new Error('useFiles must be within a FileProvider');
    }

    return context;
}

interface FileProviderProps {
    children: React.ReactNode;
}

export const FileProvider: React.FC<FileProviderProps> = ({ children }) => {
    // to hold files and folders
    const [files, setFiles] = useState<{ [key: string]: FileType | FolderType }>({});
    // to hold the current active file path
    const [activeFile, setActiveFile] = useState('');
    // to hold the data of the active file
    const [activeFileData, setActiveFileData] = useState('');

    // sort files when new file/folder is added
    useEffect(() => {
        const compareFile = (x: string, y: string): number => {
            const type1 = files[x].type;
            const type2 = files[y].type;

            if (type1 === 'folder' && type2 === 'folder') {
                return -1;
            } else if (type1 !== 'folder' && type2 === 'folder') {
                return 1;
            } else {
                return x.localeCompare(y);
            }
        }

        const sortedFileKeys = Object.keys(files).sort(compareFile);

        const sortedFiles = sortedFileKeys.reduce((acc, key) => {
            acc[key] = files[key];
            return acc;
        },
            {} as { [key: string]: FileType | FolderType }
        );
        setFiles(sortedFiles);
    }, [Object.keys(files).length]);

    // add files/folders to the state
    const setFilesFn = (dir: string, tempFileNames: string[]) => {
        tempFileNames.forEach((file) => {
            const fileData = {
                type: file.includes('.') ? 'file' : 'folder',
                path: `${dir ? dir + '\\' : ''}${file}`,
                data: ''
            };

            setFiles((prevFiles) => {
                return {
                    ...prevFiles,
                    [fileData.path]: fileData
                }
            });
        });
    }

    // to update data of a file and set it active
    const setFileData = (dir: string, data: string) => {
        setActiveFile(dir);
        setActiveFileData(data);
        setFiles((prevFiles) => {
            return {
                ...prevFiles,
                [dir]: {
                    ...prevFiles[dir],
                    data
                }
            }
        });
    }

    // remove file from state
    const removeFile = (dir: string) => {
        setFiles((prevFiles) => {
            const newFiles = { ...prevFiles };
            delete newFiles[dir];
            return newFiles;
        });
    }

    // remove folder and its data from state
    const removeFolder = (dir: string) => {
        setFiles((prevFiles) => {
            const newFiles = { ...prevFiles };
            Object.keys(newFiles).forEach((file) => {
                if (file.startsWith(dir)) {
                    delete newFiles[file];
                }
            });
            return newFiles;
        });
    }

    return (
        <FileUtils.Provider
            value={{
                files,
                setFiles: setFilesFn,
                setFileData,
                activeFile,
                activeFileData,
                removeFile,
                removeFolder
            }}
        >
            {children}
        </FileUtils.Provider>
    );
}
