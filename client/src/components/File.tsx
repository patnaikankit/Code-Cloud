import { useState } from "react"

interface FileProps {
    path: string
}

const File: React.FC<FileProps> = ({path}) => {
    const [folderHover, setFolderHover] = useState(false)

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
    }

    return(
        <li className={`text-white flex justify-between cursor-pointer relative`}
        onMouseEnter={() => {
            setFolderHover(true)
        }}
        onMouseLeave={() => {
            setFolderHover(false)
        }}>
            <p
                onClick={() => {
                    
                }}
            ></p>
        </li>
    );
}