import { useState, useEffect } from 'react'
import QRCode from 'qrcode'
import { 
  Download, 
  Copy, 
  ExternalLink, 
  Maximize2, 
  X, 
  CheckCircle,
  Image,
  Settings2
} from 'lucide-react'

export interface QRPreviewProps {
  url: string
  size?: number
  showActions?: boolean
  showUrl?: boolean
  className?: string
  format?: 'png' | 'svg'
  backgroundColor?: string
  foregroundColor?: string
  margin?: number
  errorCorrectionLevel?: 'L' | 'M' | 'Q' | 'H'
}

const QRPreview = ({
  url,
  size = 200,
  showActions = true,
  showUrl = true,
  className = '',
  format = 'png',
  backgroundColor = '#ffffff',
  foregroundColor = '#000000',
  margin = 4,
  errorCorrectionLevel = 'M'
}: QRPreviewProps) => {
  const [qrDataUrl, setQrDataUrl] = useState<string>('')
  const [svgString, setSvgString] = useState<string>('')
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string>('')
  const [copied, setCopied] = useState(false)
  const [isModalOpen, setIsModalOpen] = useState(false)

  const generateQR = async () => {
    if (!url.trim()) {
      setError('URL is required')
      setIsLoading(false)
      return
    }

    setIsLoading(true)
    setError('')

    try {
      const options = {
        width: size,
        margin: margin,
        color: {
          dark: foregroundColor,
          light: backgroundColor
        },
        errorCorrectionLevel: errorCorrectionLevel
      }

      if (format === 'png') {
        const dataUrl = await QRCode.toDataURL(url, options)
        setQrDataUrl(dataUrl)
        setSvgString('')
      } else {
        const svg = await QRCode.toString(url, {
          ...options,
          type: 'svg'
        })
        setSvgString(svg)
        setQrDataUrl('')
      }
    } catch (err) {
      console.error('QR Code generation failed:', err)
      setError('Failed to generate QR code')
    } finally {
      setIsLoading(false)
    }
  }

  const downloadQR = () => {
    if (!qrDataUrl && !svgString) return

    const link = document.createElement('a')
    const filename = `qr-code-${Date.now()}`
    
    if (format === 'png' && qrDataUrl) {
      link.href = qrDataUrl
      link.download = `${filename}.png`
    } else if (format === 'svg' && svgString) {
      const blob = new Blob([svgString], { type: 'image/svg+xml' })
      const svgUrl = URL.createObjectURL(blob)
      link.href = svgUrl
      link.download = `${filename}.svg`
    }
    
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    
    // Clean up blob URL for SVG
    if (format === 'svg') {
      setTimeout(() => URL.revokeObjectURL(link.href), 100)
    }
  }

  const copyQRImage = async () => {
    if (!qrDataUrl) return

    try {
      // Convert data URL to blob
      const response = await fetch(qrDataUrl)
      const blob = await response.blob()
      
      // Copy to clipboard if browser supports it
      if (navigator.clipboard && window.ClipboardItem) {
        await navigator.clipboard.write([
          new ClipboardItem({
            [blob.type]: blob
          })
        ])
        setCopied(true)
        setTimeout(() => setCopied(false), 2000)
      } else {
        // Fallback: download the image
        downloadQR()
      }
    } catch (err) {
      console.error('Failed to copy QR code:', err)
      // Fallback to download
      downloadQR()
    }
  }

  const openInNewTab = () => {
    if (qrDataUrl) {
      const newWindow = window.open()
      if (newWindow) {
        newWindow.document.write(`
          <html>
            <head><title>QR Code</title></head>
            <body style="margin:0;padding:20px;text-align:center;background:#f5f5f5;">
              <img src="${qrDataUrl}" alt="QR Code" style="max-width:100%;height:auto;" />
            </body>
          </html>
        `)
      }
    } else if (svgString) {
      const newWindow = window.open()
      if (newWindow) {
        newWindow.document.write(`
          <html>
            <head><title>QR Code</title></head>
            <body style="margin:0;padding:20px;text-align:center;background:#f5f5f5;">
              ${svgString}
            </body>
          </html>
        `)
      }
    }
  }

  // Generate QR code when component mounts or props change
  useEffect(() => {
    generateQR()
  }, [url, size, format, backgroundColor, foregroundColor, margin, errorCorrectionLevel])

  return (
    <>
      <div className={`bg-white border border-gray-200 rounded-lg overflow-hidden ${className}`}>
        {/* QR Code Display */}
        <div className="relative">
          <div 
            className="flex items-center justify-center p-4 bg-gray-50 cursor-pointer hover:bg-gray-100 transition-colors"
            onClick={() => setIsModalOpen(true)}
          >
            {isLoading ? (
              <div 
                className="flex items-center justify-center"
                style={{ width: size, height: size }}
              >
                <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
              </div>
            ) : error ? (
              <div 
                className="flex items-center justify-center text-red-500 text-sm"
                style={{ width: size, height: size }}
              >
                <div className="text-center">
                  <Settings2 className="h-8 w-8 mx-auto mb-2" />
                  <p>{error}</p>
                </div>
              </div>
            ) : (
              <div className="relative group">
                {format === 'png' && qrDataUrl ? (
                  <img
                    src={qrDataUrl}
                    alt="QR Code"
                    className="max-w-full h-auto"
                    style={{ width: size, height: size }}
                  />
                ) : format === 'svg' && svgString ? (
                  <div
                    dangerouslySetInnerHTML={{ __html: svgString }}
                    style={{ width: size, height: size }}
                  />
                ) : null}
                
                {/* Overlay on hover */}
                <div className="absolute inset-0 bg-black bg-opacity-0 group-hover:bg-opacity-10 transition-all duration-200 flex items-center justify-center">
                  <Maximize2 className="h-6 w-6 text-white opacity-0 group-hover:opacity-100 transition-opacity" />
                </div>
              </div>
            )}
          </div>

          {/* Actions */}
          {showActions && !isLoading && !error && (qrDataUrl || svgString) && (
            <div className="absolute top-2 right-2 flex space-x-1">
              <button
                onClick={(e) => {
                  e.stopPropagation()
                  downloadQR()
                }}
                className="p-2 bg-white bg-opacity-90 hover:bg-opacity-100 text-gray-700 rounded-md shadow-sm border border-gray-200 transition-all"
                title="Download QR Code"
              >
                <Download className="h-4 w-4" />
              </button>
              
              {format === 'png' && (
                <button
                  onClick={(e) => {
                    e.stopPropagation()
                    copyQRImage()
                  }}
                  className={`p-2 bg-white bg-opacity-90 hover:bg-opacity-100 rounded-md shadow-sm border border-gray-200 transition-all ${
                    copied ? 'text-green-600' : 'text-gray-700'
                  }`}
                  title="Copy QR Code"
                >
                  {copied ? <CheckCircle className="h-4 w-4" /> : <Copy className="h-4 w-4" />}
                </button>
              )}
              
              <button
                onClick={(e) => {
                  e.stopPropagation()
                  openInNewTab()
                }}
                className="p-2 bg-white bg-opacity-90 hover:bg-opacity-100 text-gray-700 rounded-md shadow-sm border border-gray-200 transition-all"
                title="Open in new tab"
              >
                <ExternalLink className="h-4 w-4" />
              </button>
            </div>
          )}
        </div>

        {/* URL Info */}
        {showUrl && (
          <div className="p-3 bg-white border-t border-gray-200">
            <div className="flex items-start space-x-2">
              <Image className="h-4 w-4 text-gray-500 mt-0.5 flex-shrink-0" />
              <div className="flex-1 min-w-0">
                <p className="text-xs font-medium text-gray-700 mb-1">Target URL:</p>
                <p className="text-xs text-gray-600 break-all">{url}</p>
              </div>
            </div>
          </div>
        )}
      </div>

      {/* Modal for enlarged view */}
      {isModalOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-75">
          <div className="relative max-w-lg w-full mx-4">
            <div className="bg-white rounded-lg overflow-hidden">
              {/* Modal Header */}
              <div className="flex items-center justify-between p-4 border-b border-gray-200">
                <h3 className="text-lg font-semibold text-gray-900">QR Code</h3>
                <button
                  onClick={() => setIsModalOpen(false)}
                  className="p-1 text-gray-400 hover:text-gray-600 rounded"
                >
                  <X className="h-5 w-5" />
                </button>
              </div>
              
              {/* Modal Content */}
              <div className="p-6 text-center">
                <div className="flex items-center justify-center mb-4">
                  {format === 'png' && qrDataUrl ? (
                    <img
                      src={qrDataUrl}
                      alt="QR Code"
                      className="max-w-full h-auto"
                      style={{ maxWidth: 400, maxHeight: 400 }}
                    />
                  ) : format === 'svg' && svgString ? (
                    <div
                      dangerouslySetInnerHTML={{ __html: svgString }}
                      style={{ maxWidth: 400, maxHeight: 400 }}
                    />
                  ) : null}
                </div>
                
                {/* Modal Actions */}
                <div className="flex items-center justify-center space-x-3">
                  <button
                    onClick={downloadQR}
                    className="flex items-center space-x-2 px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition-colors"
                  >
                    <Download className="h-4 w-4" />
                    <span>Download {format.toUpperCase()}</span>
                  </button>
                  
                  {format === 'png' && (
                    <button
                      onClick={copyQRImage}
                      className={`flex items-center space-x-2 px-4 py-2 rounded-md transition-colors ${
                        copied 
                          ? 'bg-green-100 text-green-700' 
                          : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                      }`}
                    >
                      {copied ? <CheckCircle className="h-4 w-4" /> : <Copy className="h-4 w-4" />}
                      <span>{copied ? 'Copied!' : 'Copy'}</span>
                    </button>
                  )}
                </div>
                
                {/* URL */}
                <div className="mt-4 p-3 bg-gray-50 rounded-md text-left">
                  <p className="text-sm font-medium text-gray-700 mb-1">Target URL:</p>
                  <p className="text-sm text-gray-600 break-all">{url}</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      )}
    </>
  )
}

export default QRPreview