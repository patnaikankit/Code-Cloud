import { useState } from "react"
import { useFiles } from "../utils/FileUtil"
import { useWebsocket } from "../utils/WSutil"
import { MdDelete } from "react-icons/md"


interface FileProps {
    path: string
}

const File: React.FC<FileProps> = ({path}) => {
    const [folderHover, setFolderHover] = useState(false)
    const { getFile, sendMessage } = useWebsocket()
    const { removeFile } = useFiles()
    const { activeFile } = useFiles()

    const deleteFile = () => {
        const cmd = {
            dir: '/app/' + path
                .split('/')
                .filter((e) => e !== '')
                .slice(0, -1)
                .join('/'),
            command: 'rm' + path.split('/').pop(),
            type: 'files',
            isFile: ''
        }
        sendMessage(JSON.stringify(cmd))
        removeFile(path)
    }


    return(
        <li className={`text-white flex justify-between cursor-pointer relative ${activeFile === path ? 'bg-gray-800' : 'bg-transparent'}`}
        onMouseEnter={() => {
            setFolderHover(true)
        }}
        onMouseLeave={() => {
            setFolderHover(false)
        }}>
            <p
                onClick={() => {
                    getFile(path)
                }}
            >{
                path
                    .split('/')
                    .filter((e) => e !== '')
                    .pop()
            }</p>
            {folderHover && (
                <MdDelete 
                className="text-white text-lg cursor-pointer absolute right-0 top-[0.2rem]"
                size={'0.8rem'}
                onClick={deleteFile}
            />
            )
            }
        </li>
    );
}

export default File;