import { useState } from "react"
import { useNavigate } from 'react-router-dom'
import Navbar from "../components/Navbar"

const Home = () => {
    const [gitLink,setGitLink] = useState('')
    const [rootDir, setRootDir] = useState('')
    // tech stack
    const [stack, setStack] = useState('nextjs')
    const [isLoading, setIsLoading] = useState(false)
    const Navigate = useNavigate()

    const Clone = async () => {
        setIsLoading(true)
        try{
            const response = await fetch(`http://localhost:5000/api/git/clone?git-link=${gitLink}&root-dir=${rootDir}&stack=${stack}`,{
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: ''
            })

            const data = await response.json()
            Navigate(`/${data.repo}`, { replace: false })
        }
        catch(err){
            console.error(err)
        }
        finally{
            setIsLoading(false)
        }
    }

    return(
        <div className="bg-black">
            <Navbar />
            <div className="w-screen h-screen flex justify-center items-center flex-col gap-2">
            <div className="flex flex-col w-1/2 max-w-[30rem] gap-2">
                <input 
                    type="text" 
                    value={gitLink}
                    onChange={(e) => setGitLink(e.target.value)}
                    disabled={isLoading}
                    placeholder="Github Link"
                    className="h-10 bg-gray-300 px-2 outline-none rounded"
                />

                <input 
                    type="text" 
                    value={rootDir}
                    onChange={(e) => setRootDir(e.target.value)}
                    disabled={isLoading}
                    placeholder="Root Directory"
                    className="h-10 bg-gray-300 px-2 outline-none rounded"
                />

                <select
                    value={stack}
                    onChange={(e) => setStack(e.target.value)}
                    disabled={isLoading}
                    className="h-10 bg-gray-300 px-2 outline-none rounded"
                >
                    <option value="React">React</option>
                    <option value="nextjs">Next.js</option>
                    <option value="django">Django</option>
                </select>
            </div>
            <button
                className="bg-blue-500 text-white px-4 py-2 rounded-md ml-2"
                disabled={isLoading}
                onClick={Clone}
            >
                {isLoading ? 'Cloning...' : 'Clone'}
            </button>
            </div>
        </div>
    );
}

export default Home;
