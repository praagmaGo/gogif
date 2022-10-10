package main

import (
//  "html/template"
//  "net/http"
  "fmt"
//  "log"
  "io/ioutil"
//  "time"
//  "strconv"
  "bytes"
  "encoding/binary"
  "math"
  "os"
)

func check(e error) { //https://gobyexample.com/reading-files
    if e != nil {
        panic(e)
    }
}

func main() {
  var octet uint8
  var pointeur uint32

  if len(os.Args)!=2 {
    fmt.Printf("Usage: gogif filename\n")
    os.Exit(0)
  }

  dat, err := ioutil.ReadFile(os.Args[1])
  check(err)
//  buf := bytes.NewReader(dat)

  fmt.Printf("0000 -------------------------------Header\n")
  if bytes.Compare(dat[0:3],[]byte("GIF"))==0 {
    fmt.Println("      format                           : GIF")
  } else {
    fmt.Println("      not GIF format : "+string(dat[0:3]))
    os.Exit(0)
  }

  if bytes.Compare(dat[3:6],[]byte("87a"))==0 {
    fmt.Println("      Version                          : 87a")
  } else if bytes.Compare(dat[3:6],[]byte("89a"))==0 {
    fmt.Println("      Version                          : 89a")
  } else {
    fmt.Println("      Version unknown: "+string(dat[3:6]))
  }
//  var mySlice = []byte{244, 244, 244, 244, 244, 244, 244, 244}
// extracting int from []byte
// http://stackoverflow.com/questions/11184336/how-to-convert-from-byte-to-int-in-go-programming
//  largeur:=binary.BigEndian.Uint16(dat[6:])
  fmt.Printf("0006 -------------------------------Logical Screen Descriptor\n")
  largeur:=binary.LittleEndian.Uint16(dat[6:])
  fmt.Printf("      Logical Screen Width    (pixels) : %d\n",largeur)
  hauteur:=binary.LittleEndian.Uint16(dat[8:])
  fmt.Printf("      Logical Screen Height   (pixels) : %d\n",hauteur)
  
  var gctf bool
//  octet=binary.LittleEndian.Uint8(dat[10:])
  if dat[10] & 0x80 == 0{
    gctf=false
    fmt.Printf("      Global Color Table               : no\n")
  }else{
    gctf=true
    fmt.Printf("      Global Color Table               : yes\n")
  }
//  if gctf {fmt.Printf("toto\n")}
  octet=(dat[10] & 0x70)>>4+1
  plz:=math.Pow(2,float64(octet))
  //  dat, err := ioutil.ReadFile("gogif.gif")
  fmt.Printf("      Color resolution                 : %d (%d) (richesse of original palette)\n",octet,int16(plz))
  if dat[10] & 0x08 == 0{
    fmt.Printf("      Global Color Table sorted        : no\n")
  }else{
    fmt.Printf("      Global Color Table sorted        : yes\n")
  }
  octet=(dat[10] & 0x07)+1
  tcz:=int16(math.Pow(2,float64(octet)))
  fmt.Printf("      Global Color Table size          : 0x%x\n",tcz)
  fmt.Printf("      Background color index           : %d",dat[11])
  if gctf==true {
    fmt.Printf("\n")
  }else{
    fmt.Printf(" (ignored, should be 0)\n")
  }
  if dat[12]==0{
    fmt.Printf("      Pixel Aspect Ratio               : no information\n")
  }else{
    fmt.Printf("      To implement: Pixel Aspect Ratio calculation\n")
  }
  pointeur=13
  if gctf==true {
    // on saute la Global Color Table
    fmt.Printf("000d ----------------------Global Color Table\n")
    fmt.Printf("      Size                             : 0x%04x\n",3*uint32(tcz))
    pointeur+=3*uint32(tcz)
  }
//  fmt.Printf("Block id           : %x\n",dat[pointeur])
  for {
    if dat[pointeur]==0x21 {
      fmt.Printf("%04x -------------------------------Extension block (0x21)\n",pointeur)
      if dat[pointeur+1]==0xf9 {
        fmt.Printf("%04x -----------------------Graphic Control Extension (0xf9)\n",pointeur+1)

        if dat[pointeur+2]!=4 {fmt.Print("**** Erreur size is not 4 ****")}

        octet=(dat[pointeur+3] & 0x1c)>>2
        fmt.Print("      Disposal                         : ")
        if octet==0 {fmt.Print("non specified")
        }else if octet==1 {fmt.Print("Do not dispose")
        }else if octet==2 {fmt.Print("Restore to background")
        }else if octet==3 {fmt.Print("Restore to previous")
        }else {fmt.Printf("To be defined (%02x)",octet)
        }
        fmt.Println("")
        
        if dat[pointeur+3] & 0x02 == 0{
          fmt.Printf("      User Input                       : no\n")
        }else{
          fmt.Printf("      User Input                       : yes\n")
        }
        if dat[pointeur+3] & 0x01 == 0{
          fmt.Printf("      Transparent Color                : no\n")
        }else{
          fmt.Printf("      Transparent Color                : yes\n")
        }
        delay:=binary.LittleEndian.Uint16(dat[pointeur+4:])
        fmt.Printf("      Delay Time                       : %d\n",delay)
        fmt.Printf("      Transparent Color Index          : %02x\n",dat[pointeur+6])
 

        pointeur+=8
//        os.Exit(0)
      }else if dat[pointeur+1]==0xfe {
        fmt.Printf("%04x -----------------------Comment Label (0xfe)\n",pointeur+1)
        for {
          size:=uint32(dat[pointeur+2])
          fmt.Printf("Taille = %d\n",size)
          pointeur+=2
          if size==0 {pointeur++;break}
          fmt.Printf("      Message                          : \""+string(dat[pointeur:pointeur+size+1])+"\"\n")
          pointeur+=size-1
        }
//        fmt.Printf("Sortie boucle for\n")
      }else if dat[pointeur+1]==0xff {
        fmt.Printf("%04x -----------------------Application extension (0xff)\n",pointeur+1)
        fmt.Printf("Application data\n")
//        pointeur+=14
        pointeur+=2
        for{
          fmt.Printf("      Bloc taille %02x\n",dat[pointeur])
          if dat[pointeur]==0 {pointeur+=1;break}
          pointeur+=uint32(dat[pointeur])+1
        }
      }else {
        fmt.Printf("Graphic Control Label           : erreur (%x) (should be 0xf9)\n",dat[pointeur+1])
      }
    }else if dat[pointeur]==0x2c {
      fmt.Printf("%04x -------------------------------Image Descriptor (0x2c)\n",pointeur)

      globol:=binary.LittleEndian.Uint16(dat[pointeur+1:])
      fmt.Printf("      Left Position           (pixels) : %d\n",globol)
      globol=binary.LittleEndian.Uint16(dat[pointeur+3:])
      fmt.Printf("      Top Position            (pixels) : %d\n",globol)
      globol=binary.LittleEndian.Uint16(dat[pointeur+5:])
      fmt.Printf("      Image Width             (pixels) : %d\n",globol)
      globol=binary.LittleEndian.Uint16(dat[pointeur+7:])
      fmt.Printf("      Image Height            (pixels) : %d\n",globol)
      
      if dat[pointeur+9] & 0x40 == 0{
        fmt.Printf("      Interlace                        : no\n")
      }else{
        fmt.Printf("      Interlace                        : yes\n")
      }
      if dat[pointeur+9] & 0x20 == 0{
        fmt.Printf("      Sorted                           : no\n")
      }else{
        fmt.Printf("      Sorted                           : yes\n")
      }

      if dat[pointeur+9] & 0x80 == 0{
        fmt.Printf("      Local Color Table present        : no\n")
        pointeur+=10
      }else{
        fmt.Printf("      Local Color Table present        : yes\n")
        pointeur+=9
        octet=(dat[pointeur] & 0x07)+1
        tcz:=int32(math.Pow(2,float64(octet)))
        fmt.Printf("      Local Color Table size           : 0x%x\n",tcz)
        pointeur++
        // on saute la Local Color Table
        fmt.Printf("%04x -----------------------Local Color Table\n",pointeur)
        fmt.Printf("      Size                             : 0x%04x\n",3*uint32(tcz))
        pointeur+=3*uint32(tcz)
      }
       
//      fmt.Printf("Pointeur  : %x\n",pointeur)
      
      fmt.Printf("%04x -----------------------Image data LZW\n",pointeur)
      fmt.Printf("      LZW code size                    : %x\n",dat[pointeur])
      pointeur++
        for{
          fmt.Printf("      Bloc taille %02x\n",dat[pointeur])
          if dat[pointeur]==0 {pointeur+=1;break}
          pointeur+=uint32(dat[pointeur])+1
        }
      
//      os.Exit(0)
    }else if dat[pointeur]==0x3b {
      fmt.Printf("%04x -------------------------------Gif Trailer (0x3b)\n",pointeur)
      os.Exit(0)
    }else if dat[pointeur]==0x02 {
      fmt.Printf("------------------------------------Image lzw (0x02)\n")
      pointeur+=1
      for{
        fmt.Printf("Bloc taille %02x\n",dat[pointeur])
        if dat[pointeur]==0 {pointeur+=1;break}
        pointeur+=uint32(dat[pointeur])+1
      }
//      os.Exit(0)
    }else if dat[pointeur]==0x00 {
      fmt.Printf("id 0\n")
      break
    }else {fmt.Printf("Unknown bloc so far (%x)\n",dat[pointeur]);os.Exit(0)}
  }
  
  fmt.Printf("------------------Stop\n")
}
