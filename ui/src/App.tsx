import './App.css'
import Designer from './components/Designer/Index'

function App() {

  return (
    <div className="App">
      <div className="masthead">
        <div className="grid">
          <div className="col-10 col:sm-12">
            <span className="title text-2xl">Greypot Studio<span className="text-yellow-500">-dev</span></span>
          </div>
          {/* <div className="col-2 col:sm-12">
            <ul className="list-none">
              <li><a href="#">Examples</a></li>
              <li><a href="#">Settings</a></li>
              <li><a href="#">About</a></li>
              <li><a href="#">Contribute on Github</a></li>
            </ul>
          </div> */}
        </div>
      </div>
      <Designer />
      <footer>
        <p>👋 Mulibwanji</p>
        <a href="https://nndi.cloud/oss/greypot">Greypot Studio</a> is an open-source project brought to you from Malawi by NNDI.
      </footer>
    </div>
  )
}

export default App
