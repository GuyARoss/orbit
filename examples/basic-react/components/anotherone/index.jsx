import React, { useState } from 'react'

const AnotherOne = () => {
    const [counter,setCounter] = useState(0)

    return (
        <>
            <div>You have pressed me {counter} times.</div>
            <button onClick={() => setCounter((prev) => prev+1)}>Click Me</button>
        </>
    )
}

export default AnotherOne