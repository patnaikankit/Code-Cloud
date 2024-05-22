// folder interface

import React, { useState } from "react"
import { useFiles } from "../utils/FileUtil"
import { useWebsocket } from "../utils/WSutil"
import { IoMdArrowDropdown, IoMdArrowDropup } from "react-icons/io"
import { BiPlus } from "react-icons/bi"
import { MdDelete } from "react-icons/md"
import { CgClose } from "react-icons/cg"
import File from "./File"

interface Folderprops{
    path: string
}

const Folder: React.FC<Folderprops> = ({ path }) => {
    const [folderHover, setFolderHover] = useState(false)
    // track if the file/folder form is clicked
    const [addClick, setAddClick] = useState(false)
    const [addHover, setAddHover] = useState(false)
    // if a file or folder is added
    const [addType, setAddType] = useState<'folder' | 'file'>('file')
    const [addName, setAddName] = useState('')
    // to check if a folder is open
    const [isOpen, setIsOpen] = useState(false)
    const { files, removeFolder } = useFiles()
    const { sendMessage } = useWebsocket()

    // to create a new file or folder
    const createFileOrFolder = (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault()

        if(addName === '' || addName.includes('/') || addName.includes(' ') || addName.includes('.')){
            return
        }

        const cmd = {
            dir: '/app/' + path,
            command: 
                addType === 'file' ? 'touch' + addName + ' && echo ' + addName : 'mkdir' + addName + ' && echo ' + addName,
            type: 'files',
            isFile: addType
        }
        sendMessage(JSON.stringify(cmd))
        setAddClick(false)
        setAddName('')
    }

    // to delete a given folder
    const deleteFolder = () => {
        const cmd = {
            dir: '/app/' + path
                            .split('/')
                            .filter((e) => e !== '')
                            .slice(0, -1)
                            .join('/'),
            command: 'rm-r' + path.split('/').pop(),
            type: 'files',
            isFile: ''
        }
        sendMessage(JSON.stringify(cmd))
        removeFolder(path)
    }

    return(
        <li
            className="text-white flex flex-col relative"
            onMouseEnter={() => {
                setFolderHover(true)
            }}
            onMouseLeave={() => {
                setFolderHover(false)
            }}
        >
            <div
                className="flex justify-between cursor-pointer"
                onClick={() => {
                    setIsOpen((e) => !e)
                    const cmd = {
                        dir: '/app/' + path,
                        command: 'ls',
                        type: 'files',
                        isFile: ''
                    }
                    sendMessage(JSON.stringify(cmd))
                }}
            >
                <p>
                    {path
                        .split('/')
                        .filter((e) => e !== '')
                        .pop()}
                </p>
                {isOpen ? <IoMdArrowDropdown /> : <IoMdArrowDropup />}
            </div>
            {folderHover && (
                <>
                <BiPlus 
                    className="text-white text-lg cursor-pointer absolute right-8" onClick={() => {
                        setAddHover((e) => !e)
                        setIsOpen(true)
                    }}
                />
                <MdDelete 
                    className="text-white text-lg cursor-pointer absolute right-4 top-[0.2rem]" size={'0.8rem'} onClick={() => {
                        deleteFolder()
                    }}
                />
                </>
            )}
            {addHover && (
                <span
                    className="absolute right-0 top-8 bg-gray-300 flex flex-col w-24 items-center z-50 text-gray-700"
                    onMouseLeave={() => {
                        setAddHover(false)
                    }}
                >
                    <p
                        className="hover:bg-gray-600 hover:text-white w-24 flex items-center justify-center cursor-pointer"
                        onClick={() => {
                            setAddType('folder')
                            setAddClick(true)
                        }}
                    >
                        Folder
                    </p>

                    <p
                        className="hover:bg-gray-600 hover:text-white w-24 flex items-center justify-center cursor-pointer"
                        onClick={() => {
                            setAddType('file')
                            setAddClick(true)
                        }}
                    >
                        File
                    </p>
                </span>
            )}

            {isOpen && (
                <ul className="ml-2">{
                    addClick && (
                        <form action="border-[0.1rem] relative text-white text-sm" onSubmit={createFileOrFolder}>
                            <input type="text"
                            className="w-full bg-transparent outline-none pl-1"
                            value={addName}
                            onChange={(e) => {
                                setAddName(e.target.value)
                            }}
                            />
                            <p className="absolute right-[1.2rem] top-[0.02rem] text-xs">
                                *{addType}
                            </p>

                            <CgClose 
                                className="absolute right-0 top-0 text-white text-lg"
                                onClick={() => {
                                    setAddClick(false)
                                }}
                            />
                        </form>
                    )
                }
                {
                    Object.keys(files).map((file) => {
                        if(!files[file].path){
                            return null
                        }
                        if(files[file].path.includes(path) && files[file].path.split('/').filter((e) => e !== '').length === path.split('/').filter((e) => e !== '').length + 1){
                            if(files[file].type === 'file'){
                                return <File key={file} path={path}/>
                            }
                            else{
                                return(
                                    <Folder 
                                        key={file}
                                        path={files[file].path.split('/app')[1] || files[file].path}
                                    />
                                )
                            }
                        }
                        return null;
                    })
                }
                </ul>
            )}
        </li>
    );
}

export default Folder;