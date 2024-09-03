- [ABA routing transit number](https://en.wikipedia.org/wiki/ABA_routing_transit_number)
    - [American Bankers Association (ABA)](https://en.wikipedia.org/wiki/American_Bankers_Association)
    - MICR Format: XXXXYYYYC
        - XXXX: Federal Reserve Routing Symbol
            - First 2 digits - Federal reserve banks
            - Electronic Transaction Identifiers (ETIs)
            - Thrift: credit unions and savings banks, etc...
        - YYYY: ABA Institution Identifier
        - C: Check Digit
         - `(3(d1 + d4 + d7) + 7(d2 + d5 + d8) + (d3 + d6 + d9)) mod 10 = 0`
    - Fraction Format: PP-YYYY/XXXX
        - PP: a 1 or 2 digit Prefix, bank's check processing center location
            - 1-49: processing center located in major city
            - 50-99: processing center located in non-major city
        - XXXX: Federal Reserve Routing Symbol
        - YYYY: ABA Institution Identifier
- [Magnetic ink character recognition (MICR)](https://en.wikipedia.org/wiki/Magnetic_ink_character_recognition)
    - Why it's called `TOAD`: The characters `T` (Transit), `O` (Onus), `A` (Amount), and `D` (Dash) are commonly found within the MICR code, hence the nickname "TOAD".
    - E-13B 
        - 14-character set
        - 10 decimal digits
        - and symbols:
            |Symbol|Unicode|Name|
            |-|-|-|
            | ⑆ | U+2446 | OCR Branch Bank Identification|
            | ⑈ | U+2447 | OCR Amount of Check |
            | ⑇ | U+2448 | OCR Dash |
            | ⑉ | U+2449 | OCR Customer Account Number |
    - CMC-7
- [Optical character recognition (OCR)](https://en.wikipedia.org/wiki/Optical_character_recognition)
    - [OCR-A](https://en.wikipedia.org/wiki/OCR-A)
    - [OCR-B](https://en.wikipedia.org/wiki/OCR-B)
- [Payment card number](https://en.wikipedia.org/wiki/Payment_card_number)
    - payment card number or primary account number (PAN)
    - found on: 
        - credit cards
        - debit cards
        - store-value cards
        - gift cards
- [Business Identifier Codes (BIC)](https://en.wikipedia.org/wiki/ISO_9362)
    - ISO 9362
- [Society for Worldwide Interbank Financial Telecommunication (Swift)](https://en.wikipedia.org/wiki/SWIFT)
- [Universal Payment Identification Code (UPIC)](https://en.wikipedia.org/wiki/Universal_Payment_Identification_Code)
- [Automated clearing house (ACH)](https://en.wikipedia.org/wiki/Automated_clearing_house)
- [Federal Reserve Wire Network (Fedwire)](https://en.wikipedia.org/wiki/Fedwire)
- [Clearing House Interbank Payments System (CHIPS)](https://en.wikipedia.org/wiki/Clearing_House_Interbank_Payments_System)
- [National Automated Clearing House (NACH)](https://en.wikipedia.org/wiki/National_Automated_Clearing_House)
- Issuer identification number (IIN)
- [Luhn algorithm](https://en.wikipedia.org/wiki/Luhn_algorithm)
