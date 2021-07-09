import React, { useEffect } from 'react'
import ReactDOM from 'react-dom';


const App = ({ test }) => {
    console.log(test)

    return (
        <>
            Is dis Working?
        </>
    )
}

ReactDOM.render(
    <App {...JSON.parse(document.getElementById('xx_orbit_data').textContent)} />,
    document.getElementById('root')
);