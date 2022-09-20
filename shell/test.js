function myFunction() {
  try {
    var sheet = SpreadsheetApp.getActiveSpreadsheet();
    var date = new Date();
    var ss = sheet.getSheetByName("JIRA_Refresh_Daily_6pm");
    let startRow = 2; // First row of data to process
    let numRows = ss.getLastRow(); // Number of rows to process
    // Fetch the range of cells A2:B30
    var all = new Map();
    var m = new Map();
    for (i = startRow; i <= numRows; i++) {
        index = "B" + i;
        m.set(ss.getRange(index).getValue(), index);
    }
    //console.log(m)
    const dataRange = ss.getRange(startRow, 1, numRows, 14);
    // Fetch values for each row in the Range.
    const data = dataRange.getValues();
    for (let a of data) {
      if (all.has(a[0])) {
         all.get(a[0]).push(a);
      } else {
        all.set(a[0], []);
        all.get(a[0]).push(a);  
      } 
    }
    all.delete("");
    month = date.getMonth() + 1;
    now = date.getDate() + '/'+ month + '/' + date.getFullYear();
    //createNewReport
    //var newSheet = SpreadsheetApp.create("report-" + now);
    var newSheet = sheet.insertSheet("TestDailyReport-" + now);
    //writeFields
    //newSheet.getRange("A0").setBackground("#efefef");
    //newSheet.getRange("B0").setBackground("#efefef");
    newSheet.getRange("A1").setValue("Email Address").setFontSize("6");
    newSheet.getRange("B1").setValue("Subject").setFontSize("6");
    newSheet.getRange("D2").setValue('DELIVERY PROJECTS DAILY REPORT - ' + now).setFontColor("#45818e").setFontSize("30");
    newSheet.getRange("D4").setValue('PROJECT DETAILS').setFontSize("9").setFontColor("WHITE").setBackground("#cccccc");
    newSheet.getRange("E4").setBackground("#cccccc");
    newSheet.getRange("F4").setBackground("#cccccc");
    newSheet.getRange("G4").setBackground("#cccccc");
    newSheet.getRange("H4").setBackground("#cccccc");
    newSheet.getRange("I4").setValue('UPDATES').setFontColor("WHITE").setFontSize("9").setBackground("#3c6e9a");
    newSheet.getRange("J4").setValue('ACTIONS').setFontColor("WHITE").setFontSize("9").setBackground("#45818e");
    ar = ["JIRA TICKET", "TASK NAME", "ASSIGNEE", "STATUS", "TODAY'S UPDATES/COMMENTS", "NEXT STEPS"];
    newSheet.getRange("E5").setValue(ar[0]).setFontSize("9").setBackground("#f3f3f3");
    newSheet.getRange("F5").setValue(ar[1]).setFontSize("9").setBackground("#f3f3f3");
    newSheet.getRange("G5").setValue(ar[2]).setFontSize("9").setBackground("#f3f3f3");
    newSheet.getRange("H5").setValue(ar[3]).setFontSize("9").setBackground("#f3f3f3");
    newSheet.getRange("I5").setValue(ar[4]).setFontSize("9").setBackground("#f3f3f3").setFontColor("#3c6e9a");
    newSheet.getRange("J5").setValue(ar[5]).setFontSize("9").setFontColor("#45818e").setBackground("#f3f3f3");
    newSheet.getRange("D5").setBackground("#f3f3f3");
    var initClientTitleOrdinateIndex = 6;
    //write  
    for (let [key, value] of all) {
      var eventLength = value.length;
      var clientTitle = key;
      var clientTitleOrdinateIndex = 'D' + initClientTitleOrdinateIndex;
      newSheet.getRange(clientTitleOrdinateIndex).setValue(clientTitle).setFontSize("16");
      var eventOrdinateIndex = initClientTitleOrdinateIndex + 1;
      var index = initClientTitleOrdinateIndex + 1;
      //write event index
      for (j = 0; j < eventLength; j++) {
        var start = 'D' + index;
        newSheet.getRange(start).setValue(j + 1);
        index++
      }  
      //write event    
      for (let e of value) {     
        var mail = "=VLOOKUP(G" + eventOrdinateIndex + ",Ref_Email_List!A:B,2,False)";   
        var jiraTicketIndex = 'E' + eventOrdinateIndex;
        var taskNameIndex = 'F' + eventOrdinateIndex;
        var assigneeIndex = 'G' + eventOrdinateIndex;
        var status = 'H' + eventOrdinateIndex;
        var commentsIndex = 'I' + eventOrdinateIndex;
        var nextStepsIndex = 'J' + eventOrdinateIndex;
        var vlookupIndex = 'A' + eventOrdinateIndex;
        newSheet.getRange(jiraTicketIndex).setFormula(getUrl(m.get(e[1]))).setFontSize("10");
        newSheet.getRange(taskNameIndex).setValue(e[2]).setFontSize("10");
        newSheet.getRange(assigneeIndex).setValue(e[3]).setFontSize("10");
        newSheet.getRange(status).setValue(e[4]).setFontSize("10");
        newSheet.getRange(commentsIndex).setValue(e[7]).setFontSize("10").setWrap(false).setWrapStrategy(SpreadsheetApp.WrapStrategy.CLIP);
        newSheet.getRange(nextStepsIndex).setValue(e[3]).setFontSize("10");
        newSheet.getRange(vlookupIndex).setValue(mail).setFontSize("6");
        eventOrdinateIndex++
      }    
      //clientOrdinateIndexCount
      initClientTitleOrdinateIndex = initClientTitleOrdinateIndex + eventLength + 1;
    }
  } catch(err) {
    Logger.log(err);
  }
}

function getUrl(a) {
  try {
    var sheet = SpreadsheetApp.getActiveSpreadsheet();
    var ss = sheet.getSheetByName("JIRA_Refresh_Daily_6pm");
    u = ss.getRange(a).getFormulas().toString().split("(").toString().split(")").toString().split(",");
    var link = "=HYPERLINK(" + u[1] + ", " + u[2]+ ")";
    return link;
  }catch(err) {
    Logger.log(err);
  }
}

