import { useState } from 'react';
import './App.css';
import MappingForm from './components/MappingForm';
import MappingTable from './components/MappingTable';

function App() {
  const [count, setCount] = useState(0)

  return (
    <div className="container mx-auto p-4">
      <h1 className="text-2xl font-bold mb-4">Path to URL Mappings</h1>
      <MappingForm />
      <MappingTable />
      <div>
        <button onClick={() => setCount((count) => count + 1)}>
          count is {count}
        </button>
      </div>
    </div>
  )
}

export default App
