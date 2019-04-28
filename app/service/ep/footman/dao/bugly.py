
from browsermobproxy import Server
from selenium import webdriver
import time
import  os
from selenium.webdriver.chrome.options import Options

# configuration
#browsermobPath = './browsermob-proxy-2.1.4/bin/browsermob-proxy'
browsermobPath = 'D:\\fyf\\tool\\browsermob-proxy-2.1.4\\bin\\browsermob-proxy'
username = '972360526'
password = '61241623FYFzwq'
tokenFile = os.getcwd()+'./token.conf'
cookiesFile = os.getcwd()+'./cookie.conf'
chromedriver = os.getcwd()+"./chromedriver.exe"

def writeResult(filePath, fileContext):
    if os.path.exists(filePath):
        os.remove(filePath)
    f = open(filePath, 'w')
    f.write(fileContext)
    print(fileContext)
    f.close()
    return

def GetCookieAndToken():
    server = Server(browsermobPath)
    server.start()
    proxy = server.create_proxy()
    profile  = webdriver.FirefoxProfile()
    profile.set_proxy(proxy.selenium_proxy())

    chrome_options = Options()

    chrome_options.add_argument('--ignore-certificate-errors')
    chrome_options.add_argument('--proxy-server={0}'.format(proxy.proxy))

    os.environ["webdriver.chrome.driver"] = chromedriver

    driver = webdriver.Chrome(chromedriver,chrome_options=chrome_options)

    #driver = webdriver.PhantomJS(firefox_profile=profile,executable_path = geckodriverPah)

    proxy.new_har("bugly", options={"captureHeaders":True})

    driver.get("https://bugly.qq.com/v2/")
    time.sleep(3)
    driver.find_element_by_class_name("login_btn").click()
    time.sleep(3)
    driver.switch_to.frame("ptlogin_iframe")
    time.sleep(3)
    driver.find_element_by_id("switcher_plogin").click()
    time.sleep(3)
    driver.find_element_by_id("u").send_keys(username)
    time.sleep(3)
    driver.find_element_by_id("p").clear()
    driver.find_element_by_id("p").send_keys(password)
    time.sleep(3)
    driver.find_element_by_id("login_button").click()
    time.sleep(10)
    driver.find_element_by_xpath('//*[@id="root"]/div/div/div[2]/div/div/div/div[2]/table/tbody/tr/td[1]/div/div[1]/img').click()
    time.sleep(3)
    driver.find_element_by_xpath('//*[@id="root"]/div/div/div[2]/div/div[1]/div[2]/ul[2]/li/a').click()
    time.sleep(10)
    strCookies = ""
    strToken = ""
    cookies = driver.get_cookies()
    requestDict = proxy.har['log']['entries']

    for index in range(len(requestDict)):
        for k in requestDict[index]:
            if k == "request" and requestDict[index][k]['url'].find('v2/issueList')>=0:
                for inn in range(len(requestDict[index][k]['headers'])):
                    for ik in requestDict[index][k]['headers'][inn]:
                        if ik == 'name' and requestDict[index][k]['headers'][inn][ik]=='X-token':
                            strToken = requestDict[index][k]['headers'][inn]['value']
                        if ik == 'name' and requestDict[index][k]['headers'][inn][ik]=='Cookie' and requestDict[index][k]['headers'][inn]['value'].find('pt2gguin')>=0 and requestDict[index][k]['headers'][inn]['value'].find('bugly_session')>=0 and requestDict[index][k]['headers'][inn]['value'].find('referrer')>=0:
                            strCookies = requestDict[index][k]['headers'][inn]['value']
                    if strToken!="" and strCookies!="":
                        break;


    writeResult(tokenFile,strToken)
    writeResult(cookiesFile,strCookies)

    server.stop()

if __name__ == '__main__':
    GetCookieAndToken()

