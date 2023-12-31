import { toBase64 } from 'fast-base64'
import { MouseEvent, useEffect, useRef, useState } from 'react'
import { TabView, TabPanel } from 'primereact/tabview'
import { Button } from 'primereact/button'
import { Message } from 'primereact/message'
import CodeMirror from '@uiw/react-codemirror';
import { html } from '@codemirror/lang-html';
import { javascript } from '@codemirror/lang-javascript';
import { json } from '@codemirror/lang-json';
import { vscodeDark } from '@uiw/codemirror-theme-vscode'
import TemplatePreview from './TemplatePreview'
import { ExampleHTML, ExampleData } from './examples'

function App() {
  const [dataCode, setDataCode] = useState(ExampleData)
  const [templateCode, setTemplateCode] = useState(ExampleHTML)
  const [editorTheme, setEditorTheme] = useState(vscodeDark)

  const downloadRef: any = useRef(null)

  const [downloadName, setDownloadName] = useState('test.pdf')

  async function downloadTemplateLocally(event: MouseEvent<HTMLButtonElement>): Promise<boolean> {
    let utf8Encode = new TextEncoder();
    let templateAsBase64 = await toBase64(utf8Encode.encode(templateCode))
    downloadRef.current.setAttribute("href", `data:application/octet-stream;base64,${templateAsBase64}`)
    downloadRef.current.setAttribute("download", 'greypot-template.html')
    return await downloadRef.current.click();
  }

  async function uploadAndRenderPDF(event: MouseEvent<HTMLButtonElement>): Promise<boolean> {
    event.preventDefault()

    const templateRequest = {
      Name: 'test.html',
      Template: templateCode,
      Data: JSON.parse(dataCode),
    }

    let response = await fetch(`/_studio/generate/pdf/${templateRequest.Name}`, {
      method: 'POST',
      mode: 'cors',
      cache: 'no-cache',
      credentials: 'same-origin',
      headers: {
        'Content-Type': 'application/json',
        'X-Greypot-Studio-Version': '0.0.1-dev',
      },
      redirect: 'error',
      referrerPolicy: 'no-referrer',
      body: JSON.stringify(templateRequest),
    });

    if (response.ok) {
      type ExportResponse = {
        data: string,
        type: string,
        reportId: string
      }
      let res = await response.json() as ExportResponse;

      downloadRef.current.setAttribute("href", `data:application/octet-stream;base64,${res.data}`)
      downloadRef.current.setAttribute("download", templateRequest.Name.replace(".html", ".pdf"))
      await downloadRef.current.click();
    }

    return false
  }

  async function uploadAndRenderExcel(event: MouseEvent<HTMLButtonElement>): Promise<boolean> {
    event.preventDefault()

    const templateRequest = {
      Name: 'test.html',
      Content: templateCode
    }

    let response = await fetch("/_studio/upload-template", {
      method: 'POST',
      mode: 'cors',
      cache: 'no-cache',
      credentials: 'same-origin',
      headers: {
        'Content-Type': 'application/json',
        'X-Greypot-Studio-Version': '0.0.1-dev',
      },
      redirect: 'error',
      referrerPolicy: 'no-referrer',
      body: JSON.stringify(templateRequest),
    });

    if (response.ok) {
      let testDataJSON = JSON.parse(dataCode)
      let response = await fetch(`/_studio/reports/export/excel/${templateRequest.Name}`, {
        method: 'POST',
        mode: 'cors',
        cache: 'no-cache',
        credentials: 'same-origin',
        headers: {
          'Content-Type': 'application/json',
          'X-Greypot-Studio-Version': '0.0.1-dev',
        },
        redirect: 'error',
        referrerPolicy: 'no-referrer',
        body: JSON.stringify(testDataJSON),
      });

      if (response.ok) {
        type ExportResponse = {
          data: string,
          type: string,
          reportId: string
        }
        let res = await response.json() as ExportResponse;

        downloadRef.current.setAttribute("href", `data:application/octet-stream;base64,${res.data}`)
        downloadRef.current.setAttribute("download", templateRequest.Name.replace(".html", ".xlsx"))
        await downloadRef.current.click();
      }
    }

    return false
  }

  return (
    <>
      <div className="grid">
        <div className="col-12 md:col-4 sm:col-12 xl:col-8">
          <TabView>
            <TabPanel header="Template Editor">
              <CodeMirror
                width="100%"
                height="400px"
                extensions={[html(), javascript()]}
                value={templateCode}
                onChange={(e) => setTemplateCode(e)}
                theme={editorTheme}
              // options={options}
              // editorDidMount={editorDidMount}
              />
            </TabPanel>
            <TabPanel header="Preview">
              <TemplatePreview htmlTemplateCode={templateCode} testData={dataCode} />
            </TabPanel>
          </TabView>
        </div>
        <div className="col-12 md:col-4 sm:col-12 xl:col-4">
          <TabView>
            <TabPanel header="Test Data">
              <CodeMirror
                width="100%"
                height="400px"
                extensions={[json()]}
                theme={editorTheme}
                value={dataCode}
                onChange={(e) => setDataCode(e)}
              />
            </TabPanel>
            <TabPanel header="cURL Request">
              <code lang='bash'>
                curl -i -H "Content-Type: application/json" -X POST "https://greypot-studio.fly.dev/_studio/reports/export/pdf/test.html" -d &nbsp;
                "{dataCode.replaceAll('"', '\\"')}"

              </code>
            </TabPanel>
            <TabPanel header="Schema">
              Coming soon...
            </TabPanel>
          </TabView>
        </div>
      </div>
      <div className="action-area p-3">
        <a style={{ display: 'none' }} ref={downloadRef} download={downloadName}></a>
        <Button label="Download Template" onClick={downloadTemplateLocally} />&nbsp;
        <Button label="PDF Preview with Test Data" onClick={uploadAndRenderPDF} />&nbsp;
        {/* <Button label="Excel Preview with Test Data" onClick={uploadAndRenderExcel} /> */}
      </div>
    </>
  )
}

export default App
