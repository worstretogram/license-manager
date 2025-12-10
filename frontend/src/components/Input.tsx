'use client'
const Input:React.FC<{value: string, placeholder: string, onChange: (text: string) => void}> = ({value, onChange, placeholder}) => {

    return (
            <input value={value} placeholder={placeholder} onChange={(e) => onChange(e.target.value)} className=" bg-secondary w-full h-9 rounded-input px-4 outline-none cursor-pointer text-white font-normal focus:border-borderActive transition duration-200" type="text"/>
    )
}

export default Input