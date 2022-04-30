import React from 'react'

export default ({ title, description, color = 'default'}) => {
    return (
        <div className={`notice-${color}`}>
            <strong>{title}</strong>
            <p>{description}</p>
        </div>
    )
}