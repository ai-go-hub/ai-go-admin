interface anyObj {
    [key: string]: any
}

interface ApiResponse<T = any> {
    code: number
    message: string
    time: number
    data: T
}

interface Window {
    loading: boolean
}

type Writeable<T> = { -readonly [P in keyof T]: T[P] }
