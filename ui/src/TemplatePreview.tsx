import { useState } from "react"
import { twig } from 'twig'

interface TemplatePreviewProps {
    htmlTemplateCode: string 
    testData: string
}

export default function({ htmlTemplateCode, testData }: TemplatePreviewProps) {
    const template = twig({
        data: htmlTemplateCode
    })

    let dataParsed = JSON.parse(testData)
    let renderedTemplate = template.render({ data: dataParsed })

    const [isFullScreen, setIsFullScreen] = useState(false)
    function toggleFullScreen() {
        setIsFullScreen(!isFullScreen)
    }

    return (
        <div className={isFullScreen ? 'preview-wrapper fullscreen' : 'preview-wrapper '}>
            <iframe
                width="100%"
                style={{ height: isFullScreen ? "90vh" : "400px" }}
                srcDoc={renderedTemplate} />
            <button onClick={toggleFullScreen}>{isFullScreen ? 'Close Full Screen' : 'Full Screen'}</button>
        </div>
    )
}