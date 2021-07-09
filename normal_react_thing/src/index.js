import React from 'react'
import ReactDOM from 'react-dom';

const App = () => {
    return (
        <>
            Is dis Working?
        </>
    )
}

console.log('loaded', App, document.getElementById('root'))

ReactDOM.render(<App />, document.getElementById('root'));