import React, { useState, useEffect} from 'react'

export default function() {
    const [packages, usePackages] = useState([])
    
    useEffect(() => {
        const getData = async () => {
            const res = await fetch('api/index')
            const json = await res.json()
            console.log(json)
            return json
        }
        const json = getData().then(json => usePackages(json))
        //usePackages(json)
    }, [])

    return (
        <div>
            Index
            <ul>
                {packages.map((pkg) => <li key={pkg}>{pkg}</li>)}
            </ul>
        </div>
    )
}