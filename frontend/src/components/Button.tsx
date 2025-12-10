'use client'
const Button:React.FC<{title?: string, onClick?: () => void}> = ({title = 'label', onClick}) => {    
    return (
        <button onClick={onClick} className="h-8 px-4 w-fit border border-border rounded-input text-white font-bold">{title}</button>
    )
}

export default Button