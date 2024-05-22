import './App.css'
import { BrowserRouter } from 'react-router-dom'
import RoutesContainer from './helpers/Routes'

function App() {
  return (
    <BrowserRouter>
      <RoutesContainer />
    </BrowserRouter>
  )
}

export default App
