import { License } from '@/types'
import { makeAutoObservable} from 'mobx'

class Store {
    constructor() {
        makeAutoObservable(this)
    }
    isAdding: boolean = false
    username: string = ""
    password: string = ""
    owner: string = ""
    maxMessages: string = ""
    maxUsers: string = ""
    licenses:License[] = []

    changeIsAdding = (status: boolean) => {
        this.isAdding = status
    }

    changeLicenses = (data: License[]) => {
        this.licenses = data
    }

    changeUsername = (text: string) => {        
        this.username = text
    }
    changePassword = (text: string) => {
        this.password = text
    }
    changeOwner = (text: string) => {
        this.owner = text
    }
    changeMaxMessages = (text: string) => {
        this.maxMessages = text
    }
    changeMaxUsers = (text: string) => {
        this.maxUsers = text
    }

}

export default new Store()