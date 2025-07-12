import { useState, useEffect } from 'react'
import QRCode from 'qrcode'
import { Download, RefreshCw, Palette, Settings2, Image } from 'lucide-react'

export interface QRCodeOptions {
  size: number
  backgroundColor: string
  foregroundColor: string
  margin: number
  errorCorrectionLevel: 'L' | 'M' | 'Q' | 'H'
  type: 'image/png' | 'image/svg+xml'
}

export interface QRGeneratorProps {
  url: string
  className?: string
  showCustomization?: boolean
  onGenerated?: (dataUrl: string, options: QRCodeOptions) => void
  initialOptions?: Partial<QRCodeOptions>
}

const defaultOptions: QRCodeOptions = {
  size: 256,
  backgroundColor: '#ffffff',
  foregroundColor: '#000000',
  margin: 4,
  errorCorrectionLevel: 'M',
  type: 'image/png'
}

const QRGenerator = ({
  url,
  className = '',
  showCustomization = true,
  onGenerated,
  initialOptions = {}
}: QRGeneratorProps) => {
  const [options, setOptions] = useState<QRCodeOptions>({
    ...defaultOptions,
    ...initialOptions
  })
  const [qrDataUrl, setQrDataUrl] = useState<string>('')
  const [svgString, setSvgString] = useState<string>('')
  const [isGenerating, setIsGenerating] = useState(false)
  const [showSettings, setShowSettings] = useState(false)
  const [error, setError] = useState<string>('')

  const generateQR = async () => {
    if (!url.trim()) {
      setError('URL is required')
      return
    }

    setIsGenerating(true)
    setError('')

    try {
      if (options.type === 'image/png') {
        // Generate PNG
        const dataUrl = await QRCode.toDataURL(url, {
          width: options.size,
          margin: options.margin,
          color: {
            dark: options.foregroundColor,
            light: options.backgroundColor
          },
          errorCorrectionLevel: options.errorCorrectionLevel
        })
        setQrDataUrl(dataUrl)
        setSvgString('')
        onGenerated?.(dataUrl, options)
      } else {
        // Generate SVG
        const svg = await QRCode.toString(url, {
          type: 'svg',
          width: options.size,
          margin: options.margin,
          color: {
            dark: options.foregroundColor,
            light: options.backgroundColor
          },
          errorCorrectionLevel: options.errorCorrectionLevel
        })
        setSvgString(svg)
        setQrDataUrl('')
        onGenerated?.(svg, options)
      }
    } catch (err) {
      console.error('QR Code generation failed:', err)
      setError('Failed to generate QR code. Please check the URL and try again.')
    } finally {
      setIsGenerating(false)
    }
  }

  const downloadQR = () => {
    if (!qrDataUrl && !svgString) return

    const link = document.createElement('a')
    
    if (options.type === 'image/png' && qrDataUrl) {
      link.href = qrDataUrl
      link.download = `qr-code-${Date.now()}.png`
    } else if (options.type === 'image/svg+xml' && svgString) {
      const blob = new Blob([svgString], { type: 'image/svg+xml' })
      const svgUrl = URL.createObjectURL(blob)
      link.href = svgUrl
      link.download = `qr-code-${Date.now()}.svg`
    }
    
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    
    // Clean up blob URL
    if (options.type === 'image/svg+xml') {
      setTimeout(() => URL.revokeObjectURL(link.href), 100)
    }
  }

  const updateOption = <K extends keyof QRCodeOptions>(
    key: K,
    value: QRCodeOptions[K]
  ) => {
    setOptions(prev => ({ ...prev, [key]: value }))
  }

  const resetToDefaults = () => {
    setOptions({ ...defaultOptions, ...initialOptions })
  }

  // Auto-generate on mount and when options change
  useEffect(() => {
    generateQR()
  }, [url, options])

  return (
    <div className={`space-y-4 ${className}`}>
      {/* Header */}
      <div className="flex items-center justify-between">
        <h3 className="text-lg font-semibold text-gray-900">QR Code Generator</h3>
        {showCustomization && (
          <div className="flex items-center space-x-2">
            <button
              onClick={() => setShowSettings(!showSettings)}
              className={`p-2 rounded-md transition-colors ${
                showSettings 
                  ? 'bg-blue-100 text-blue-600' 
                  : 'bg-gray-100 hover:bg-gray-200 text-gray-600'
              }`}
              title="Customize QR Code"
            >
              <Settings2 className="h-4 w-4" />
            </button>
            <button
              onClick={resetToDefaults}
              className="p-2 bg-gray-100 hover:bg-gray-200 text-gray-600 rounded-md transition-colors"
              title="Reset to defaults"
            >
              <RefreshCw className="h-4 w-4" />
            </button>
          </div>
        )}
      </div>

      {/* Error Message */}
      {error && (
        <div className="p-3 bg-red-50 border border-red-200 rounded-md text-red-700 text-sm">
          {error}
        </div>
      )}

      {/* Customization Panel */}
      {showCustomization && showSettings && (
        <div className="bg-gray-50 border border-gray-200 rounded-lg p-4 space-y-4">
          <h4 className="font-medium text-gray-900 flex items-center">
            <Palette className="h-4 w-4 mr-2" />
            Customization Options
          </h4>
          
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            {/* Size */}
            <div>
              <label htmlFor="qr-size" className="block text-sm font-medium text-gray-700 mb-1">
                Size (px)
              </label>
              <input
                id="qr-size"
                type="range"
                min="128"
                max="512"
                step="32"
                value={options.size}
                onChange={(e) => updateOption('size', parseInt(e.target.value))}
                className="w-full"
              />
              <div className="text-xs text-gray-500 mt-1">{options.size}px</div>
            </div>

            {/* Margin */}
            <div>
              <label htmlFor="qr-margin" className="block text-sm font-medium text-gray-700 mb-1">
                Margin
              </label>
              <input
                id="qr-margin"
                type="range"
                min="0"
                max="8"
                step="1"
                value={options.margin}
                onChange={(e) => updateOption('margin', parseInt(e.target.value))}
                className="w-full"
              />
              <div className="text-xs text-gray-500 mt-1">{options.margin} units</div>
            </div>

            {/* Foreground Color */}
            <div>
              <label htmlFor="qr-fg-color" className="block text-sm font-medium text-gray-700 mb-1">
                Foreground Color
              </label>
              <div className="flex items-center space-x-2">
                <input
                  id="qr-fg-color"
                  type="color"
                  value={options.foregroundColor}
                  onChange={(e) => updateOption('foregroundColor', e.target.value)}
                  className="w-12 h-8 border border-gray-300 rounded cursor-pointer"
                />
                <input
                  type="text"
                  value={options.foregroundColor}
                  onChange={(e) => updateOption('foregroundColor', e.target.value)}
                  className="flex-1 px-3 py-1 border border-gray-300 rounded-md text-sm"
                  placeholder="#000000"
                />
              </div>
            </div>

            {/* Background Color */}
            <div>
              <label htmlFor="qr-bg-color" className="block text-sm font-medium text-gray-700 mb-1">
                Background Color
              </label>
              <div className="flex items-center space-x-2">
                <input
                  id="qr-bg-color"
                  type="color"
                  value={options.backgroundColor}
                  onChange={(e) => updateOption('backgroundColor', e.target.value)}
                  className="w-12 h-8 border border-gray-300 rounded cursor-pointer"
                />
                <input
                  type="text"
                  value={options.backgroundColor}
                  onChange={(e) => updateOption('backgroundColor', e.target.value)}
                  className="flex-1 px-3 py-1 border border-gray-300 rounded-md text-sm"
                  placeholder="#ffffff"
                />
              </div>
            </div>

            {/* Error Correction Level */}
            <div>
              <label htmlFor="qr-error-level" className="block text-sm font-medium text-gray-700 mb-1">
                Error Correction
              </label>
              <select
                id="qr-error-level"
                value={options.errorCorrectionLevel}
                onChange={(e) => updateOption('errorCorrectionLevel', e.target.value as 'L' | 'M' | 'Q' | 'H')}
                className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm"
              >
                <option value="L">Low (7%)</option>
                <option value="M">Medium (15%)</option>
                <option value="Q">Quartile (25%)</option>
                <option value="H">High (30%)</option>
              </select>
            </div>

            {/* Format */}
            <div>
              <label htmlFor="qr-format" className="block text-sm font-medium text-gray-700 mb-1">
                Format
              </label>
              <select
                id="qr-format"
                value={options.type}
                onChange={(e) => updateOption('type', e.target.value as 'image/png' | 'image/svg+xml')}
                className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm"
              >
                <option value="image/png">PNG (Raster)</option>
                <option value="image/svg+xml">SVG (Vector)</option>
              </select>
            </div>
          </div>
        </div>
      )}

      {/* QR Code Preview */}
      <div className="flex flex-col items-center space-y-4">
        <div className="bg-white border-2 border-gray-200 rounded-lg p-4 shadow-sm">
          {isGenerating ? (
            <div className="flex items-center justify-center" style={{ width: options.size, height: options.size }}>
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
            </div>
          ) : (
            <div className="flex items-center justify-center">
              {options.type === 'image/png' && qrDataUrl ? (
                <img
                  src={qrDataUrl}
                  alt="QR Code"
                  className="max-w-full h-auto"
                  style={{ width: options.size, height: options.size }}
                />
              ) : options.type === 'image/svg+xml' && svgString ? (
                <div
                  dangerouslySetInnerHTML={{ __html: svgString }}
                  style={{ width: options.size, height: options.size }}
                />
              ) : null}
            </div>
          )}
        </div>

        {/* Actions */}
        {(qrDataUrl || svgString) && !isGenerating && (
          <div className="flex items-center space-x-3">
            <button
              onClick={downloadQR}
              className="flex items-center space-x-2 px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition-colors"
            >
              <Download className="h-4 w-4" />
              <span>Download {options.type === 'image/png' ? 'PNG' : 'SVG'}</span>
            </button>
            
            <button
              onClick={generateQR}
              disabled={isGenerating}
              className="flex items-center space-x-2 px-4 py-2 bg-gray-100 text-gray-700 rounded-md hover:bg-gray-200 transition-colors disabled:opacity-50"
            >
              <RefreshCw className={`h-4 w-4 ${isGenerating ? 'animate-spin' : ''}`} />
              <span>Regenerate</span>
            </button>
          </div>
        )}
      </div>

      {/* URL Info */}
      <div className="bg-gray-50 border border-gray-200 rounded-md p-3">
        <div className="flex items-start space-x-2">
          <Image className="h-4 w-4 text-gray-500 mt-0.5 flex-shrink-0" />
          <div className="flex-1 min-w-0">
            <p className="text-sm font-medium text-gray-700">Target URL:</p>
            <p className="text-sm text-gray-600 break-all">{url}</p>
          </div>
        </div>
      </div>
    </div>
  )
}

export default QRGenerator