import Foundation
import AppKit
import SceneKit
import Metal
import AVFoundation
import CoreText
import ImageIO
import UniformTypeIdentifiers

// Deterministic, genuinely lit and extruded 3D brand animation.
// render-brand-motion <opening|closing> <silent.mp4> <font-directory> [--preview <directory>]
let args = CommandLine.arguments
guard args.count >= 4, ["opening", "closing"].contains(args[1]) else {
    fatalError("Usage: render-brand-motion <opening|closing> <output.mp4> <fonts> [--preview <directory>]")
}
let opening = args[1] == "opening"
let mode = args[1]
let outputURL = URL(fileURLWithPath: args[2])
let previewOnly = args.count == 6 && args[4] == "--preview"
let previewURL = previewOnly ? URL(fileURLWithPath: args[5]) : outputURL.deletingLastPathComponent()
let width = 1920, height = 1080
let fps: Int32 = 30
let duration = opening ? 9.0 : 8.0
for filename in ["manrope-700.ttf", "dm-sans-400.ttf", "dm-sans-700.ttf"] {
    CTFontManagerRegisterFontsForURL(URL(fileURLWithPath: args[3]).appendingPathComponent(filename) as CFURL, .process, nil)
}
func color(_ hex: UInt32, _ alpha: CGFloat = 1) -> NSColor {
    NSColor(srgbRed: CGFloat((hex >> 16) & 255)/255, green: CGFloat((hex >> 8) & 255)/255,
            blue: CGFloat(hex & 255)/255, alpha: alpha)
}
let ink = color(0x071D23), white = color(0xF7FCF8), amber = color(0xF0AD4E), muted = color(0xA4C4C1)
func material(_ hex: UInt32, shine: CGFloat = 0.3) -> SCNMaterial {
    let m = SCNMaterial()
    m.lightingModel = .blinn
    m.diffuse.contents = color(hex)
    m.specular.contents = color(0xBDDAD2)
    m.shininess = shine
    return m
}
let tealMat = material(0x12483F, shine: 0.65)
let paperMat = material(0xF7FCF8, shine: 0.2)
let foldMat = material(0xABCFC2)
let amberMat = material(0xF0AD4E, shine: 0.55)
let lineMat = material(0x195B4D)
let scene = SCNScene()

// A quiet radial studio background, retaining enough contrast for titles.
let bg = NSImage(size: NSSize(width: width, height: height))
bg.lockFocus()
ink.setFill(); NSRect(x: 0, y: 0, width: width, height: height).fill()
let gradient = NSGradient(colors: [color(0x184941), color(0x0C2A30), ink])!
gradient.draw(fromCenter: NSPoint(x: 600, y: 540), radius: 0,
              toCenter: NSPoint(x: 650, y: 520), radius: 1200, options: [.drawsAfterEndingLocation])
bg.unlockFocus()
scene.background.contents = bg
scene.lightingEnvironment.intensity = 0.4

func box(_ w: CGFloat, _ h: CGFloat, _ d: CGFloat, _ bevel: CGFloat, _ mat: SCNMaterial) -> SCNNode {
    let geometry = SCNBox(width: w, height: h, length: d, chamferRadius: bevel)
    geometry.chamferSegmentCount = 12
    geometry.materials = [mat]
    return SCNNode(geometry: geometry)
}
func vector(_ x: Double, _ y: Double, _ z: Double) -> SCNVector3 { SCNVector3(x,y,z) }
func shape(_ points: [(Double,Double)], depth: CGFloat, mat: SCNMaterial) -> SCNNode {
    let path = NSBezierPath()
    for (index,p) in points.enumerated() {
        let point = NSPoint(x: (p.0 - 32)/16, y: (32 - p.1)/16)
        if index == 0 { path.move(to: point) } else { path.line(to: point) }
    }
    path.close()
    let geo = SCNShape(path: path, extrusionDepth: depth)
    geo.chamferRadius = 0.016
    geo.materials = [mat]
    return SCNNode(geometry: geo)
}
func stroke(_ a: (Double,Double), _ b: (Double,Double), radius: CGFloat, z: Double, mat: SCNMaterial, parent: SCNNode) {
    let ax = (a.0-32)/16, ay = (32-a.1)/16, bx = (b.0-32)/16, by = (32-b.1)/16
    let cylinder = SCNCylinder(radius: radius, height: hypot(bx-ax,by-ay))
    cylinder.radialSegmentCount = 24; cylinder.materials = [mat]
    let node = SCNNode(geometry: cylinder)
    node.position = vector((ax+bx)/2,(ay+by)/2,z)
    node.eulerAngles.z = CGFloat(atan2(by-ay,bx-ax) - .pi/2)
    parent.addChildNode(node)
    for (x,y) in [(ax,ay),(bx,by)] {
        let ball = SCNSphere(radius: radius); ball.segmentCount = 20; ball.materials = [mat]
        let end = SCNNode(geometry: ball); end.position = vector(x,y,z); parent.addChildNode(end)
    }
}
let hero = SCNNode(); scene.rootNode.addChildNode(hero)
let tilePath = NSBezierPath(roundedRect:NSRect(x:-2,y:-2,width:4,height:4),xRadius:0.80,yRadius:0.80)
tilePath.flatness = 0.002
let tileGeometry = SCNShape(path:tilePath,extrusionDepth:0.42)
tileGeometry.chamferRadius = 0.06; tileGeometry.materials = [tealMat]
let tile = SCNNode(geometry:tileGeometry)
// Center the extrusion regardless of the geometry's local origin.
let tileBounds = tileGeometry.boundingBox
tile.position.z = -(tileBounds.min.z + tileBounds.max.z)/2
hero.addChildNode(tile)
let document = SCNNode(); hero.addChildNode(document)
let paper = shape([(18,11),(38,11),(49,22),(49,53),(18,53),(15,50),(15,14)], depth: 0.09, mat: paperMat)
paper.position.z = 0.30; document.addChildNode(paper)
let fold = shape([(38,11),(38,19),(41,22),(49,22)], depth: 0.05, mat: foldMat)
fold.position.z = 0.39; document.addChildNode(fold)
stroke((23,29),(36,29),radius: 0.095,z: 0.39,mat:lineMat,parent:document)
stroke((23,36),(31,36),radius: 0.095,z: 0.39,mat:lineMat,parent:document)
let badge = SCNNode(); hero.addChildNode(badge)
let badgeRimGeo = SCNCylinder(radius: 0.85,height: 0.14)
badgeRimGeo.radialSegmentCount = 96; badgeRimGeo.materials = [tealMat]
let badgeRim = SCNNode(geometry: badgeRimGeo)
badgeRim.eulerAngles.x = .pi/2; badgeRim.position = vector(13.0/16,-14.0/16,0.45); badge.addChildNode(badgeRim)
let badgeGeo = SCNCylinder(radius: 0.74,height: 0.13)
badgeGeo.radialSegmentCount = 96; badgeGeo.materials = [amberMat]
let badgeFace = SCNNode(geometry: badgeGeo)
badgeFace.eulerAngles.x = .pi/2; badgeFace.position = vector(13.0/16,-14.0/16,0.55); badge.addChildNode(badgeFace)
stroke((38.5,46),(43,50.5),radius:0.115,z:0.66,mat:tealMat,parent:badge)
stroke((43,50.5),(51,41.5),radius:0.115,z:0.66,mat:tealMat,parent:badge)

// Evidence cards pass behind the emblem as the opening assembles.
var cards: [SCNNode] = []
for i in 0..<3 {
    let card = box(1.08,1.42,0.10,0.045,material(i == 1 ? 0xA9CFC2 : 0xF1F7F0))
    for j in 0..<3 {
        let line = box(j == 2 ? 0.38 : 0.70,0.045,0.015,0.007,lineMat)
        line.position = vector(j == 2 ? -0.16 : 0,0.31-Double(j)*0.20,0.066)
        card.addChildNode(line)
    }
    scene.rootNode.addChildNode(card); cards.append(card)
}
var rings: [SCNNode] = []
for i in 0..<2 {
    let torus = SCNTorus(ringRadius: CGFloat(3.05+Double(i)*0.42),pipeRadius: 0.010)
    torus.ringSegmentCount = 160
    let m = material(i == 0 ? 0x568C80 : 0xB78B43)
    m.emission.contents = color(i == 0 ? 0x23473F : 0x40381F)
    torus.materials = [m]
    let node = SCNNode(geometry: torus); node.eulerAngles.x = .pi/2
    node.opacity = 0.34; scene.rootNode.addChildNode(node); rings.append(node)
}
var sparks: [SCNNode] = []
for i in 0..<28 {
    let ball = SCNSphere(radius: i % 5 == 0 ? 0.021 : 0.010)
    ball.segmentCount = 8
    let m = SCNMaterial(); m.lightingModel = .constant
    m.diffuse.contents = color(i % 4 == 0 ? 0xF0AD4E : 0x7FAFA2,0.6)
    ball.materials = [m]
    let node = SCNNode(geometry: ball); scene.rootNode.addChildNode(node); sparks.append(node)
}
func light(_ type: SCNLight.LightType, _ intensity: CGFloat, _ hex: UInt32, _ position: SCNVector3) -> SCNNode {
    let node = SCNNode(); node.light = SCNLight(); node.light!.type = type
    node.light!.intensity = intensity; node.light!.color = color(hex)
    node.position = position; scene.rootNode.addChildNode(node); return node
}
let key = light(.omni,1000,0xF3FFF9,vector(-4,6,8))
key.light!.attenuationEndDistance = 35
_ = light(.omni,500,0x9CE9D9,vector(5,1,5))
_ = light(.omni,900,0xF1BC70,vector(-3,-4,-1))
_ = light(.ambient,260,0xBFCEC8,vector(0,0,0))

let camera = SCNNode(); camera.camera = SCNCamera()
camera.camera!.usesOrthographicProjection = true
camera.camera!.orthographicScale = 4.5
camera.camera!.zNear = 0.1; camera.camera!.zFar = 80
camera.camera!.wantsHDR = false; camera.camera!.wantsExposureAdaptation = false
camera.position = vector(0,0,18); scene.rootNode.addChildNode(camera)
guard let device = MTLCreateSystemDefaultDevice() else { fatalError("Metal is required") }
let renderer = SCNRenderer(device: device, options: nil)
renderer.scene = scene; renderer.pointOfView = camera; renderer.autoenablesDefaultLighting = false
func clamp(_ t: Double) -> Double { min(1,max(0,t)) }
func p(_ t: Double,_ a: Double,_ b: Double) -> Double { clamp((t-a)/(b-a)) }
func smooth(_ t: Double) -> Double { t*t*(3-2*t) }
func out(_ t: Double) -> Double { 1-pow(1-t,3) }
func mix(_ a: Double,_ b: Double,_ t: Double) -> Double { a+(b-a)*t }
func update(_ t: Double) {
    SCNTransaction.begin(); SCNTransaction.disableActions = true
    let reveal = opening ? out(p(t,0.25,1.65)) : out(p(t,0.10,1.25))
    let settle = opening ? smooth(p(t,2.65,4.25)) : 1.0
    hero.position = vector(mix(0,-3.48,settle),0.20+0.04*sin(t*0.8),0)
    let heroScale = mix(1.25,0.88,settle) * (0.78+0.22*reveal)
    hero.scale = vector(heroScale,heroScale,heroScale); hero.opacity = CGFloat(reveal)
    hero.eulerAngles = vector(mix(-0.20,0.06,settle)+0.028*sin(t*0.7),
                             mix(-0.7,-0.12,settle)*(1-reveal*0.15)+0.035*sin(t*0.65),
                             mix(-0.10,0.025,settle))
    if !opening {
        hero.position.y += CGFloat((1-reveal)*0.30)
        hero.eulerAngles.y -= CGFloat((1-reveal)*0.55)
        hero.eulerAngles.z -= CGFloat((1-reveal)*0.12)
    }
    let paperIn = opening ? out(p(t,0.7,2.0)) : reveal
    document.position.z = CGFloat((1-paperIn)*2.5); document.opacity = CGFloat(paperIn)
    let badgeIn = opening ? out(p(t,1.35,2.55)) : out(p(t,0.4,1.6))
    badge.position = vector((1-badgeIn)*1.1,-(1-badgeIn)*0.8,(1-badgeIn)*3.2)
    badge.scale = vector(0.6+0.4*badgeIn,0.6+0.4*badgeIn,0.6+0.4*badgeIn)
    badge.opacity = CGFloat(badgeIn)
    for (i,card) in cards.enumerated() {
        let n = Double(i)
        let angle = t*0.20+n*2.1+0.4
        let radius = 3.3+(1-settle)*0.5
        card.position = vector(cos(angle)*radius-3.48*settle,sin(angle)*2.0, -1.5-n*0.5)
        card.eulerAngles = vector(-0.12,0.35*sin(angle),0.13*cos(angle))
        card.opacity = CGFloat(opening ? 0.65*smooth(p(t,0.25+n*0.2,1.2+n*0.2))*(1-settle) : 0)
    }
    for (i,ring) in rings.enumerated() {
        ring.position = vector(-3.48*settle,0.20,-2.4-Double(i)*0.25)
        ring.eulerAngles = vector(.pi/2+0.1*sin(t*0.3+Double(i)),0.12*cos(t*0.25),0)
        ring.opacity = CGFloat(0.27*reveal)
        let s = mix(1.05,0.87,settle); ring.scale = vector(s,s,s)
    }
    for (i,spark) in sparks.enumerated() {
        let n = Double(i), a = n*2.399+t*0.018, r = 3.5+Double(i%7)*0.65
        spark.position = vector(cos(a)*r-2.4*settle,sin(a)*r*0.65,-3-n.truncatingRemainder(dividingBy:4))
        spark.opacity = CGFloat((0.25+0.18*sin(n+t*0.7))*reveal)
    }
    SCNTransaction.commit()
}
func font(_ size: CGFloat, bold: Bool = false, brand: Bool = false) -> NSFont {
    let name = brand ? "Manrope-Bold" : (bold ? "DMSans-Bold" : "DMSans-Regular")
    return NSFont(name:name,size:size) ?? NSFont.systemFont(ofSize:size,weight:bold ? .bold : .regular)
}
@discardableResult
func text(_ value: String,_ x: CGFloat,_ top: CGFloat,_ size: CGFloat,_ foreground: NSColor,
          _ alpha: CGFloat,_ ctx: CGContext,bold: Bool = false,brand: Bool = false,kern: CGFloat = 0) -> CGFloat {
    let string = NSAttributedString(string:value,attributes:[.font:font(size,bold:bold,brand:brand),
        .foregroundColor:foreground.withAlphaComponent(alpha),.kern:kern])
    NSGraphicsContext.saveGraphicsState()
    NSGraphicsContext.current = NSGraphicsContext(cgContext:ctx,flipped:false)
    string.draw(at:NSPoint(x:x,y:CGFloat(height)-top-string.size().height))
    NSGraphicsContext.restoreGraphicsState()
    return string.size().width
}
func line(_ x: CGFloat,_ y: CGFloat,_ w: CGFloat,_ c: NSColor,_ ctx: CGContext) {
    ctx.setFillColor(c.cgColor); ctx.fill(CGRect(x:x,y:CGFloat(height)-y,width:w,height:1))
}
func buffer() -> CVPixelBuffer {
    var result: CVPixelBuffer?
    let attrs: [String:Any] = [kCVPixelBufferCGImageCompatibilityKey as String:true,
                               kCVPixelBufferCGBitmapContextCompatibilityKey as String:true]
    guard CVPixelBufferCreate(kCFAllocatorDefault,width,height,kCVPixelFormatType_32BGRA,attrs as CFDictionary,&result) == kCVReturnSuccess,
          let result else { fatalError("Pixel allocation failed") }; return result
}
func render(_ t: Double,_ pixel: CVPixelBuffer) -> CGImage {
    update(t)
    let shot = renderer.snapshot(atTime:t,with:CGSize(width:width,height:height),antialiasingMode:.multisampling4X)
    var rect = NSRect(x:0,y:0,width:width,height:height)
    guard let image = shot.cgImage(forProposedRect:&rect,context:nil,hints:nil) else { fatalError("Snapshot failed") }
    CVPixelBufferLockBaseAddress(pixel,[]); defer { CVPixelBufferUnlockBaseAddress(pixel,[]) }
    guard let ctx = CGContext(data:CVPixelBufferGetBaseAddress(pixel),width:width,height:height,bitsPerComponent:8,
        bytesPerRow:CVPixelBufferGetBytesPerRow(pixel),space:CGColorSpaceCreateDeviceRGB(),
        bitmapInfo:CGImageAlphaInfo.premultipliedFirst.rawValue|CGBitmapInfo.byteOrder32Little.rawValue) else { fatalError("Context failed") }
    ctx.draw(image,in:CGRect(x:0,y:0,width:width,height:height))
    let header = CGFloat(smooth(p(t,0.6,1.5)))
    text("LAPORSIAP",112,80,23,white,header,ctx,bold:true,kern:3)
    text(opening ? "EVIDENCE  /  CLARITY" : "PREPARE  /  REVIEW  /  REPORT",1245,80,20,muted,header,ctx,kern:2)
    line(112,135,1696,muted.withAlphaComponent(0.15*header),ctx)
    let title = CGFloat(out(p(t,opening ? 3.9 : 1.2,opening ? 4.8 : 2.1)))
    let x: CGFloat = 940+24*(1-title)
    text(opening ? "MAKE THE DETAILS CLEAR" : "YOUR NEXT STEP",x,320,21,amber,title,ctx,bold:true,kern:2.5)
    let brandWidth = text("Lapor",x-4,380,108,white,title,ctx,bold:true,brand:true,kern:-3)
    text("Siap",x-4+brandWidth,380,108,amber,title,ctx,bold:true,brand:true,kern:-3)
    let sub = CGFloat(smooth(p(t,opening ? 4.5 : 1.85,opening ? 5.4 : 2.7)))
    text(opening ? "Clarity starts with your evidence." : "Prepare. Review. Report.",x,545,34,white,sub,ctx)
    text(opening ? "Organize the facts. Review every detail." : "Continue through official reporting channels.",x,605,24,muted,sub,ctx)
    let footer = CGFloat(smooth(p(t,opening ? 5.2 : 2.4,opening ? 6.1 : 3.2)))
    line(x,700,700,muted.withAlphaComponent(0.22*footer),ctx)
    for (i,item) in ["01  CAPTURE","02  REVIEW","03  EXPORT"].enumerated() {
        let active = CGFloat(smooth(p(t,opening ? 5.4+Double(i)*0.23 : 2.5+Double(i)*0.23,
                                         opening ? 6.1+Double(i)*0.23 : 3.2+Double(i)*0.23)))
        text(item,x+CGFloat(i)*240,736,20,i == 1 ? amber : muted,active,ctx,bold:true,kern:1)
    }
    text(opening ? "FROM MESSAGES TO REVIEWED REPORT MATERIALS" : "THANK YOU FOR WATCHING",112,970,18,muted,footer,ctx,kern:2)
    let fade = max(1-smooth(p(t,0,0.4)),smooth(p(t,duration-0.6,duration-1.0/Double(fps))))
    if fade > 0 { ctx.setFillColor(ink.withAlphaComponent(CGFloat(fade)).cgColor); ctx.fill(CGRect(x:0,y:0,width:width,height:height)) }
    guard let result = ctx.makeImage() else { fatalError("Composite failed") }; return result
}
func png(_ image: CGImage,_ url: URL) {
    guard let dest = CGImageDestinationCreateWithURL(url as CFURL,UTType.png.identifier as CFString,1,nil) else { fatalError("PNG failed") }
    CGImageDestinationAddImage(dest,image,nil); guard CGImageDestinationFinalize(dest) else { fatalError("PNG write failed") }
}
try FileManager.default.createDirectory(at:previewURL,withIntermediateDirectories:true)
if previewOnly {
    for t in (opening ? [0.9,2.0,3.5,6.5] : [0.8,2.3,4.5,6.5]) {
        autoreleasepool { png(render(t,buffer()),previewURL.appendingPathComponent("\(mode)-\(t).png")) }
        print("Preview \(mode) \(t)"); fflush(stdout)
    }
    exit(0)
}
guard !FileManager.default.fileExists(atPath:outputURL.path) else { fatalError("Output exists; use a fresh path") }
try FileManager.default.createDirectory(at:outputURL.deletingLastPathComponent(),withIntermediateDirectories:true)
let writer = try AVAssetWriter(outputURL:outputURL,fileType:.mp4); writer.shouldOptimizeForNetworkUse = true
let settings: [String:Any] = [AVVideoCodecKey:AVVideoCodecType.h264,AVVideoWidthKey:width,AVVideoHeightKey:height,
    AVVideoCompressionPropertiesKey:[AVVideoAverageBitRateKey:16_000_000,AVVideoExpectedSourceFrameRateKey:fps,
        AVVideoMaxKeyFrameIntervalKey:30,AVVideoProfileLevelKey:AVVideoProfileLevelH264HighAutoLevel]]
let input = AVAssetWriterInput(mediaType:.video,outputSettings:settings); input.expectsMediaDataInRealTime = false
let adaptor = AVAssetWriterInputPixelBufferAdaptor(assetWriterInput:input,sourcePixelBufferAttributes:[
    kCVPixelBufferPixelFormatTypeKey as String:kCVPixelFormatType_32BGRA,
    kCVPixelBufferWidthKey as String:width,kCVPixelBufferHeightKey as String:height])
writer.add(input); guard writer.startWriting() else { fatalError("Writer failed: \(String(describing:writer.error))") }
writer.startSession(atSourceTime:.zero)
let count = Int(duration*Double(fps))
for frame in 0..<count {
    autoreleasepool {
        let pixel = buffer(), image = render(Double(frame)/Double(fps),pixel)
        if frame == 195 { png(image,outputURL.deletingLastPathComponent().appendingPathComponent("\(mode)-poster.png")) }
        while !input.isReadyForMoreMediaData {
            if writer.status == .failed { fatalError("Writer failed: \(String(describing:writer.error))") }
            Thread.sleep(forTimeInterval:0.005)
        }
        guard adaptor.append(pixel,withPresentationTime:CMTime(value:Int64(frame),timescale:fps)) else { fatalError("Append failed") }
    }
    if frame % 30 == 0 { print("\(mode): \(frame)/\(count) frames"); fflush(stdout) }
}
input.markAsFinished(); let done = DispatchSemaphore(value:0)
writer.finishWriting { done.signal() }; done.wait()
guard writer.status == .completed else { fatalError("Incomplete: \(String(describing:writer.error))") }
print("Completed \(outputURL.path)")
