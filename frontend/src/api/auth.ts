import store from "@/store/store";
import axios from "axios"

export const login = async (router: any) => {
    try {
        const res = await axios.post('http://127.0.0.1:8080/api/auth/login', {
            username: store.username,
            password: store.password
        })
        localStorage.setItem('token', res.data.token)
        
            router.push('/dashboard')
        
    } catch (error) {
        console.log(error);
        
    }
}