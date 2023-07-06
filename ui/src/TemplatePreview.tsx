import { ErrorInfo, useRef, useState } from "react"
import { twig } from 'twig'

interface TemplatePreviewProps {
    htmlTemplateCode: string 
    testData: string
}

export default function({ htmlTemplateCode, testData }: TemplatePreviewProps) {
    const [isFullScreen, setIsFullScreen] = useState(false)
    const renderResult = useRef({ isError: false, message: '' })

    const template = twig({
        data: htmlTemplateCode
    })

    let renderedTemplate = "";
    try {
        let dataParsed = JSON.parse(testData)
        renderedTemplate = template.render({ data: dataParsed })
        renderResult.current.isError = false
    } catch (e) {
        console.error("failed to render template", e)
        renderResult.current.isError = true
        renderResult.current.message = `${e}`
    }

    function toggleFullScreen() {
        setIsFullScreen(!isFullScreen)
    }

    return (
        <div className={isFullScreen ? 'preview-wrapper fullscreen' : 'preview-wrapper '}>
            { renderResult.current.isError ? 
                <div>
                    <p>There was an error processing the data or rendering the template. Please check that the data is valid JSON</p>
                    <p style={{ color: "red" }}>{renderResult.current.message}</p>
                </div>
            :  <>
                <iframe
                    width="100%"
                    style={{ height: isFullScreen ? "90vh" : "400px" }}
                    srcDoc={renderedTemplate} />
                <button onClick={toggleFullScreen}>{isFullScreen ? 'Close Full Screen' : 'Full Screen'}</button>
                </>
            }
        </div>
    )
}