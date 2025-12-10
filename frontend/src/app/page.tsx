'use client'

import { login } from "@/api/auth";
import Button from "@/components/Button";
import Input from "@/components/Input";
import store from "@/store/store";
import { observer } from "mobx-react-lite";
import { useRouter } from "next/navigation";


const Home = observer(() => {

  const router = useRouter()

  return (
    <div className="max-w-96 w-full h-fit min-h-60 flex flex-col items-center gap-4 bg-primary rounded-primary p-8">
      <h1 className="text-white text-xl font-bold">Welcome Back</h1>
     <Input placeholder="Username" value={store.username} onChange={store.changeUsername}/>
     <Input placeholder="Password" value={store.password} onChange={store.changePassword}/>
     <Button title="Login" onClick={() => login(router)}/>
    </div>
  );
})

export default Home