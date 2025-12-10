'use client'
import store from "@/store/store"
import Button from "./Button"
import Input from "./Input"
import { observer } from "mobx-react-lite"
import { createLicense } from "@/api/licence"

const AddModal = observer(() => {
  const handleBackgroundClick = () => {
    store.changeIsAdding(false)
  }

  const handleModalClick = (e: React.MouseEvent<HTMLDivElement>) => {
    e.stopPropagation()
  }

  return (
    <div
      onClick={handleBackgroundClick}
      className="fixed inset-0 bg-black/50 backdrop-blur-sm flex items-center justify-center z-50"
    >
      <div
        onClick={handleModalClick}
        className="max-w-96 w-full h-fit min-h-60 flex flex-col items-center gap-4 bg-primary rounded-primary p-8"
      >
        <h1 className="text-white text-xl font-bold">Creating a license</h1>
        <Input value={store.owner} onChange={store.changeOwner} placeholder="Owner" />
        <Input value={store.maxMessages} onChange={store.changeMaxMessages} placeholder="Max Messages"/>
        <Input value={store.maxUsers} onChange={store.changeMaxUsers} placeholder="Max Users" />
        <Button onClick={createLicense} title="Create" />
      </div>
    </div>
  )
})

export default AddModal
