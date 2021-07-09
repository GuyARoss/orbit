import React from 'react'

export default ({ name, age }) => (
    <>
        This is a test component

        <div>
            You are: {age}
        </div>
        <div>
            Your name is {name}
        </div>
    </>
)