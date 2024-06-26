import { useParams } from "react-router-dom";
import { useFiles } from "../utils/FileUtil";
import { useWebsocket } from "../utils/WSutil";
import { useEffect, useState } from "react";
import FolderStructure from "../components/FolderStructure";
import { BiSave, BiShareAlt } from "react-icons/bi";
import CodeEditor from '@uiw/react-textarea-code-editor' 
import Terminal from "../components/Terminal";

const Editor = () => {
    const { id } = useParams()
    const { activeFileData, activeFile } = useFiles()
    const { setSocket, saveFile } = useWebsocket()
    const [fileData, setFileDData] = useState(activeFileData)
    const [time, setTime] = useState<ReturnType<typeof setTimeout>  | null>(null)
    const [isSaved, setIsSaved] = useState(true)

    useEffect(() => {
        setSocket(id!)
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [id])

    return(
        <div className="w-screen h-screen bg-background flex">
            <FolderStructure />
            <div className="w-[calc(100vw-16rem)] max-h-screen ml-64">
                <div className="bg-fStructBackground w-[100%] flex justify-between p-2 pr-3 h-10">
                    <p></p>
                    <BiShareAlt 
                        color="white"
                        size={'1.2rem'}
                        className="cursor-pointer"
                        onClick={() => {
                            window.open(`http://${id}.localhost:5000/`, '_blank')
                        }}
                    />

                    <BiSave 
                        color="white"
                        size={'1.2rem'}
                        className="cursor-pointer"
                        onClick={() => {
                            saveFile(activeFile, fileData)
                            clearTimeout(time!)
                            setIsSaved(true)
                        }}
                    />
                </div>

                <div>
                    <p className="text-white px-2 text-xs">
                        {activeFile}
                        {isSaved ? '' : '*'}
                    </p>
                    <div className="overflow-auto h-[calc(100vh-19.5rem)]">
                    <CodeEditor 
                        value={activeFileData}
                        language={activeFile.split('.').pop()}
                        placeholder="Choose a File"
                        onChange={(e) => {
                            setIsSaved(false)
                            time && clearTimeout(time)
                            setFileDData(e.target.value)
                            const temp = setTimeout(() => {
                                saveFile(activeFile, e.target.value)
                                setIsSaved(true)
                            }, 4000);
                            setTime(temp)
                        }}
                        padding={15}
                        className="w-[calc(100vw-16rem)] text-white p-3 text-sm rounded-md outline-none overflow-scroll"
                        style={{
                            fontFamily: 'ui-monospace,SFMono-Regular,SF Mono,Consolas,Liberation Mono,Menlo,monospace'
                        }}
                    />
                    </div>
                    <Terminal />
                </div>
            </div>
        </div>
    );
}

export default Editor;