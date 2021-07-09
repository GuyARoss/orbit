import React from 'react'

import Test from '../components/test'

const Hello = ({ name, age }) => {
    return (
        <>
            Hello, this thing working?

            <hr />
            <Test name={name} age={age} />
        </>
    )
}

export default Hello