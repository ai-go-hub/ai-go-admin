import '@/styles/loading.scss'

export const loading = {
    show: () => {
        const div = document.createElement('div')
        div.className = 'ai-go-page-loading'
        div.innerHTML = `
            <div class="container">
                <div class="main">
                    <div class="item"></div>
                    <div class="item"></div>
                    <div class="item"></div>
                    <div class="item"></div>
                    <div class="item"></div>
                    <div class="item"></div>
                    <div class="item"></div>
                    <div class="item"></div>
                    <div class="item"></div>
                </div>
            </div>
        `
        document.body.insertBefore(div, document.body.childNodes[0])
    },
    hide: () => {
        const el = document.querySelector('.ai-go-page-loading')
        el && el.parentNode?.removeChild(el)
    },
}
